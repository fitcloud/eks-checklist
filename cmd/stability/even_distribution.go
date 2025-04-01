package stability

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckPodDistributionAndAffinity checks whether pods are evenly distributed via affinity or topologySpreadConstraints.
func CheckPodDistributionAndAffinity(clientset kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "동일한 역할을 하는 Pod를 다수의 노드에 분산 배포",
		Manual:     false,
		Passed:     true,
		FailureMsg: "일부 Pod에 affinity나 유효한 topologySpreadConstraints 설정이 누락되어 있습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	for _, pod := range pods.Items {
		affinityExists := pod.Spec.Affinity != nil
		topologyValid := false

		if len(pod.Spec.TopologySpreadConstraints) > 0 {
			topologyValid = true
			for _, constraint := range pod.Spec.TopologySpreadConstraints {
				if constraint.MaxSkew > 1 {
					topologyValid = false
					result.Resources = append(result.Resources,
						fmt.Sprintf("Namespace: %s | Pod: %s - maxSkew 값이 %d (1 초과)", pod.Namespace, pod.Name, constraint.MaxSkew))
				}
			}
		}

		if !affinityExists && !topologyValid {
			result.Resources = append(result.Resources,
				fmt.Sprintf("Namespace: %s | Pod: %s - affinity와 유효한 topologySpreadConstraints 설정이 모두 없음", pod.Namespace, pod.Name))
		}
	}

	if len(result.Resources) > 0 {
		result.Passed = false
	} else {
		result.Passed = true
	}

	return result
}

package reliability

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// SingletonPodCheck checks for standalone pods that are not managed by a controller (Deployment, StatefulSet, etc.).
func SingletonPodCheck(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "[REL-001] 싱글톤 Pod 미사용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "Deployment나 StatefulSet 등의 컨트롤러에 속하지 않은 Standalone Pod가 존재합니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/runbook/reliability/REL-001",
	}

	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	standaloneFound := false

	for _, pod := range pods.Items {
		if len(pod.OwnerReferences) == 0 {
			standaloneFound = true
			result.Resources = append(result.Resources,
				fmt.Sprintf("Namespace: %s | Pod: %s (Standalone Pod)", pod.Namespace, pod.Name))
		}
	}

	if standaloneFound {
		result.Passed = false
	} else {
		result.Passed = true
	}

	return result
}

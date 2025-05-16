package scalability

import (
	"context"
	"strings"

	"eks-checklist/cmd/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckSpotNodeTerminationHandler checks whether the Spot Termination Handler is deployed.
func CheckSpotNodeTerminationHandler(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "[SCL-003] Spot 노드 사용시 Spot 중지 핸들러 적용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "Spot Termination Handler 관련 파드를 찾을 수 없습니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/runbook/scalability/SCL-003",
	}

	pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	found := false
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, "termination-handler") {
			found = true
			result.Resources = append(result.Resources,
				"Namespace: "+pod.Namespace+" | Pod: "+pod.Name)
		}
	}

	if found {
		result.Passed = true
	} else {
		result.Passed = false
	}

	return result
}

// 변경 전 코드

// package network

// import (
// 	"context"

// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // namespace에 label을 확인하여 Readiness gate를 활성화 유무를 판단
// func CheckReadinessGateEnabled(client kubernetes.Interface) bool {
// 	// EKS 클러스터의 네임스페이스 가져오기
// 	namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	// 네임스페이스의 라벨 확인
// 	for _, namespace := range namespaces.Items {
// 		if namespace.Labels["elbv2.k8s.aws/pod-readiness-gate-inject"] == "enabled" {
// 			return true
// 		}
// 	}

// 	return false
// }

// 변경 후 코드
package network

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckReadinessGateEnabled checks if any namespace has pod readiness gate enabled.
func CheckReadinessGateEnabled(controller_installed common.CheckResult, client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName: "Pod Readiness Gate 적용",
		Manual:    false,
		Passed:    true,
		Runbook:   "https://your.runbook.url/latest-tag-image",
	}

	if !controller_installed.Passed {
		result.Passed = false
		result.FailureMsg = "AWS Load Balancer Controller가 설치되어 있지 않습니다"
		return result
	}

	namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	found := false
	for _, ns := range namespaces.Items {
		if ns.Labels["elbv2.k8s.aws/pod-readiness-gate-inject"] == "enabled" {
			result.Resources = append(result.Resources,
				fmt.Sprintf("Namespace: %s | Label: elbv2.k8s.aws/pod-readiness-gate-inject=enabled", ns.Name))
			found = true
		}
	}

	if found {
		result.Passed = true
		result.SuccessMsg = "일부 네임스페이스에 Pod Readiness Gate가 적용되어 있습니다."
	} else {
		result.Passed = false
		result.FailureMsg = "Pod Readiness Gate가 적용된 네임스페이스가 없습니다."
	}

	return result
}

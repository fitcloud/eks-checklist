// 기존 코드

// package scalability

// import (
// 	"context"
// 	"strings"

// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // Spot Node Termination Handler 설치 여부 확인
// func CheckSpotNodeTerminationHandler(client kubernetes.Interface) bool {
// 	pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
// 	if err != nil {
// 		return false
// 	}

// 	// Pod 목록에서 "termination-handler"가 포함된 이름의 파드 찾기
// 	for _, pod := range pods.Items {
// 		if strings.Contains(pod.Name, "termination-handler") {
// 			return true
// 		}
// 	}

// 	// 해당 파드가 없으면 false 반환
// 	return false
// }

// 변경 후 코드
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
		CheckName:  "Spot 노드 사용시 Spot 중지 핸들러 적용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "Spot Termination Handler 관련 파드를 찾을 수 없습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
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

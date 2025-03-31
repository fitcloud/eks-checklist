// 기존 코드

// package stability

// import (
// 	"context"
// 	"fmt"

// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // 클러스터에 Cluster Autoscaler가 설치되어 있는지 확인
// func SingletonPodCheck(client kubernetes.Interface) bool {
// 	// kube-system 네임스페이스의 모든 Deployment 목록 가져오기

// 	result := true

// 	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	for _, pod := range pods.Items {
// 		if pod.OwnerReferences == nil || len(pod.OwnerReferences) == 0 {
// 			fmt.Printf("[WARNING] Standalone Pod %s in namespace %s detected\n", pod.Name, pod.Namespace)
// 			result = true
// 		}
// 	}

// 	return result
// }

// // Standalone Pod 탐지 (Deployment, StatefulSet 등에 속하지 않은 Pod 찾기)
// func checkStandalonePods(client kubernetes.Interface) bool {
// 	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	result := false

// 	for _, pod := range pods.Items {
// 		if pod.OwnerReferences == nil || len(pod.OwnerReferences) == 0 {
// 			fmt.Printf("[WARNING] Standalone Pod %s in namespace %s detected\n", pod.Name, pod.Namespace)
// 			result = true
// 		}
// 	}

// 	return result
// }

// 변경 후 코드
package stability

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
		CheckName:  "싱글톤 Pod 미사용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "Deployment나 StatefulSet 등의 컨트롤러에 속하지 않은 Standalone Pod가 존재합니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
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

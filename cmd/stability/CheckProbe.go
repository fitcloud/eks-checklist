// 변경 전 코드

// package stability

// import (
// 	"context"
// 	"fmt"

// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // CheckProbe - 모든 Pod을 검색하여 startupProbe, livenessProbe, readinessProbe 가 모두 설정되었는지 확인
// func CheckProbe(client kubernetes.Interface) bool {
// 	// 모든 네임스페이스의 Pod을 조회 (kube-system 제외는 필터로 직접 처리)
// 	// 왜냐하면 시스템 애드온 파드들은 기본적으로 몇개씩 없음 어떻게 할지는 추후 더 고민

// 	//기존 코드에서 수정 이유 : fake 클라이언트는 필드셀렉터가 안먹어서 일관된 동작체크를 하기 위함
// 	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
// 	if err != nil {
// 		fmt.Println("Error retrieving pods:", err)
// 		return false
// 	}

// 	allProbesSet := true

// 	// 각 Pod 확인
// 	for _, pod := range pods.Items {
// 		if pod.Namespace == "kube-system" {
// 			continue // kube-system 네임스페이스 파드는 검사에서 제외
// 		}

// 		for _, container := range pod.Spec.Containers {
// 			missingProbes := []string{}

// 			if container.StartupProbe == nil {
// 				missingProbes = append(missingProbes, "startupProbe")
// 			}
// 			if container.LivenessProbe == nil {
// 				missingProbes = append(missingProbes, "livenessProbe")
// 			}
// 			if container.ReadinessProbe == nil {
// 				missingProbes = append(missingProbes, "readinessProbe")
// 			}

// 			if len(missingProbes) > 0 {
// 				fmt.Printf("Pod %s in namespace %s is missing the following probes: %v\n", pod.Name, pod.Namespace, missingProbes)
// 				allProbesSet = false
// 			}
// 		}
// 	}

// 	return allProbesSet
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

// CheckProbe - 모든 Pod을 검색하여 startupProbe, livenessProbe, readinessProbe 가 모두 설정되었는지 확인
func CheckProbe(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "Probe(Startup, Readiness, Liveness) 적용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "일부 컨테이너에 startup/liveness/readiness probe가 누락되어 있습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	for _, pod := range pods.Items {
		if pod.Namespace == "kube-system" {
			continue // kube-system은 검사 제외
		}

		for _, container := range pod.Spec.Containers {
			var missing []string
			if container.StartupProbe == nil {
				missing = append(missing, "startupProbe")
			}
			if container.LivenessProbe == nil {
				missing = append(missing, "livenessProbe")
			}
			if container.ReadinessProbe == nil {
				missing = append(missing, "readinessProbe")
			}

			if len(missing) > 0 {
				result.Passed = false
				result.Resources = append(result.Resources,
					fmt.Sprintf("Namespace: %s | Pod: %s | Container: %s (미설정: %v)", pod.Namespace, pod.Name, container.Name, missing))
			}
		}
	}

	return result
}

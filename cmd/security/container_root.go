// 변경 이전 코드
// package security

// import (
// 	"context"
// 	"fmt"
// 	"strings"

// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // CheckContainerExecutionUser checks if any container is running as root (UID 0).
// func CheckContainerExecutionUser(client kubernetes.Interface) bool {
// 	// 검사에서 제외할 문자열들을 지정
// 	excludeStrings := []string{
// 		"aws-node",
// 		"coredns",
// 		"eks-pod-identity-agent",
// 		"kube-proxy",
// 	}

// 	// 모든 네임스페이스에서 파드를 조회
// 	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	// 모든 네임스페이스에서 파드를 리스트
// 	for _, pod := range pods.Items {
// 		// Pod 이름에 excludeStrings 배열의 문자열이 포함된 경우 건너뜁니다.
// 		exclude := false
// 		for _, excludeString := range excludeStrings {
// 			if strings.Contains(pod.Name, excludeString) {
// 				exclude = true
// 				break
// 			}
// 		}

// 		if exclude {
// 			continue
// 		}

// 		for _, container := range pod.Spec.Containers {
// 			if container.SecurityContext != nil && container.SecurityContext.RunAsUser != nil {
// 				if *container.SecurityContext.RunAsUser == 0 {
// 					fmt.Printf("⚠️  Root user detected in Pod: %s, Container: %s\n", pod.Name, container.Name)
// 					return true
// 				} else if container.SecurityContext.WindowsOptions != nil && container.SecurityContext.WindowsOptions.RunAsUserName != nil {
// 					// RunAsUserName이 nil이 아닌지 먼저 확인한 후 값 비교
// 					if *container.SecurityContext.WindowsOptions.RunAsUserName == "Administrator" {
// 						// Windows 환경에서 "Administrator" 사용자로 실행 중인 경우
// 						fmt.Printf("⚠️  Root user detected in Pod: %s, Container: %s\n", pod.Name, container.Name)
// 						return true
// 					}
// 				}
// 			} else {
// 				// RunAsUser가 명시되지 않은 경우, 컨테이너는 기본적으로 루트로 실행될 가능성이 있음
// 				fmt.Printf("⚠️  RunAsUser not set in Pod: %s, Container: %s (Possibly running as root)\n", pod.Name, container.Name)
// 			}
// 		}
// 	}

// 	return false
// }

// 변경 이후 코드

package security

import (
	"context"
	"fmt"
	"strings"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckContainerExecutionUser checks if any container is running as root (UID 0).
func CheckContainerExecutionUser(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "루트 유저가 아닌 유저로 컨테이너 실행",
		Manual:     false,
		Passed:     true,
		FailureMsg: "일부 컨테이너가 root 유저로 실행 중이거나, RunAsUser가 명시되지 않았습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	// 검사에서 제외할 문자열
	excludeStrings := []string{
		"aws-node",
		"coredns",
		"eks-pod-identity-agent",
		"kube-proxy",
	}

	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	for _, pod := range pods.Items {
		exclude := false
		for _, excludeString := range excludeStrings {
			if strings.Contains(pod.Name, excludeString) {
				exclude = true
				break
			}
		}

		if exclude {
			continue
		}

		for _, container := range pod.Spec.Containers {
			if container.SecurityContext != nil && container.SecurityContext.RunAsUser != nil {
				if *container.SecurityContext.RunAsUser == 0 {
					result.Passed = false
					resource := fmt.Sprintf("Namespace: %s | Pod: %s | Container: %s (명시적 root 계정 실행)", pod.Namespace, pod.Name, container.Name)
					result.Resources = append(result.Resources, resource)
				} else if container.SecurityContext.WindowsOptions != nil && container.SecurityContext.WindowsOptions.RunAsUserName != nil {
					if *container.SecurityContext.WindowsOptions.RunAsUserName == "Administrator" {
						result.Passed = false
						resource := fmt.Sprintf("Namespace: %s | Pod: %s | Container: %s (Windows Administrator 실행)", pod.Namespace, pod.Name, container.Name)
						result.Resources = append(result.Resources, resource)
					}
				}
			} else {
				result.Passed = false
				resource := fmt.Sprintf("Namespace: %s | Pod: %s | Container: %s (RunAsUser 미설정, root로 실행 가능성 존재)", pod.Namespace, pod.Name, container.Name)
				result.Resources = append(result.Resources, resource)
			}
		}
	}

	return result
}

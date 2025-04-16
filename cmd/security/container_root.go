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
		// pods label에 k8s-app 키가 있는 경우에도 패스
		if _, exists := pod.Labels["k8s-app"]; exists {
			exclude = true
		}

		// pods label에 app.kubernetes.io/managed-by 키가 있는 경우에도 패스
		if _, exists := pod.Labels["app.kubernetes.io/managed-by"]; exists {
			exclude = true
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

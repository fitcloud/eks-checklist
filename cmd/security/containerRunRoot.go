package security

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckContainerExecutionUser checks if any container is running as root (UID 0).
func CheckContainerExecutionUser(client kubernetes.Interface) bool {
	// 모든 네임스페이스에서 파드를 조회
	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 모든 네임스페이스에서 파드를 리스트
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			if container.SecurityContext != nil && container.SecurityContext.RunAsUser != nil {

				if *container.SecurityContext.RunAsUser == 0 {

					fmt.Printf("⚠️  Root user detected in Pod: %s, Container: %s\n", pod.Name, container.Name)
					return true

				} else if *container.SecurityContext.WindowsOptions.RunAsUserName == "Administrator" {

					//윈도우즈 용
					fmt.Printf("⚠️  Root user detected in Pod: %s, Container: %s\n", pod.Name, container.Name)
					return true
				}

			} else {
				// RunAsUser가 명시되지 않은 경우, 컨테이너는 기본적으로 루트로 실행될 가능성이 있음
				fmt.Printf("⚠️  RunAsUser not set in Pod: %s, Container: %s (Possibly running as root)\n", pod.Name, container.Name)

			}
		}
	}

	return false
}

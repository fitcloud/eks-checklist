package general

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	Red    = "\033[31m" // 빨간색
	Green  = "\033[32m" // 초록색
	Yellow = "\033[33m" // 노란색
	Reset  = "\033[0m"  // 기본 색상으로 리셋
)

func CheckImageTag(client kubernetes.Interface) bool {
	// 모든 네임스페이스에서 파드를 리스트
	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 모든 파드의 컨테이너 이미지를 확인
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			// 이미지에 "latest"가 포함되어 있는지 확인
			if strings.Contains(container.Image, "latest") {
				fmt.Println(Red + "✖ FAIL: Latest Tag on Container Image Found" + Reset)

				return true
			}
		}
	}

	// "latest" 이미지가 없는 경우
	fmt.Println(Green + "✔ PASS: Latest Tag on Container Image Not Found" + Reset)

	return false
}

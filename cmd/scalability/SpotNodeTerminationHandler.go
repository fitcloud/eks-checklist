package scalability

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Spot Node Termination Handler 설치 여부 확인
func CheckSpotNodeTerminationHandler(client kubernetes.Interface) bool {
	pods, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false
	}

	// Pod 목록에서 "termination-handler"가 포함된 이름의 파드 찾기
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, "termination-handler") {
			return true
		}
	}

	// 해당 파드가 없으면 false 반환
	return false
}

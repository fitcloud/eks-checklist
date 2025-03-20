package network

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// namespace에 label을 확인하여 Readiness gate를 활성화 유무를 판단
func CheckReadinessGateEnabled(client kubernetes.Interface) bool {
	// EKS 클러스터의 네임스페이스 가져오기
	namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 네임스페이스의 라벨 확인
	for _, namespace := range namespaces.Items {
		if namespace.Labels["elbv2.k8s.aws/pod-readiness-gate-inject"] == "enabled" {
			return true
		}
	}

	return false
}

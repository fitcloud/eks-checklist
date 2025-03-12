package scalability

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 클러스터의 노드 인스턴스 유형이 다양한지 확인
func CheckInstanceTypes(client kubernetes.Interface) bool {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 인스턴스 유형을 저장할 맵
	instanceTypes := make(map[string]bool)

	// 각 노드의 레이블에서 인스턴스 유형을 확인
	for _, node := range nodes.Items {
		// 'beta.kubernetes.io/instance-type' 레이블을 확인
		if instanceType, exists := node.Labels["beta.kubernetes.io/instance-type"]; exists {
			instanceTypes[instanceType] = true
		}
	}

	// 다양한 인스턴스를 사용하고 있는지 확인
	return len(instanceTypes) > 1
}

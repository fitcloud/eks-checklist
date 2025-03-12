package scalability

import (
	"context"
	"fmt"
	"strings"

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

		// Fargate 작업인 경우, 인스턴스 유형으로 "fargate"를 추가
		// ProviderID가 "fargate"인 경우를 확인
		if node.Spec.ProviderID != "" && strings.Contains(node.Spec.ProviderID, "fargate") {
			instanceTypes["fargate"] = true
		}
	}

	// 인스턴스 유형들을 쉼표로 구분하여 출력
	var types []string
	for instanceType := range instanceTypes {
		types = append(types, instanceType)
	}
	fmt.Printf("Instance Types used in the cluster: %s\n", strings.Join(types, ", "))

	// 다양한 인스턴스를 사용하고 있는지 확인
	return len(instanceTypes) > 1
}

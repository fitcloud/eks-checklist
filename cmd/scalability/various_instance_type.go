// 기존 코드

// package scalability

// import (
// 	"context"
// 	"fmt"
// 	"strings"

// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // 클러스터의 노드 인스턴스 유형이 다양한지 확인
// func CheckInstanceTypes(client kubernetes.Interface) bool {
// 	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	// 인스턴스 유형을 저장할 맵
// 	instanceTypes := make(map[string]bool)

// 	// 각 노드의 레이블에서 인스턴스 유형을 확인
// 	for _, node := range nodes.Items {
// 		// 'beta.kubernetes.io/instance-type' 레이블을 확인
// 		if instanceType, exists := node.Labels["beta.kubernetes.io/instance-type"]; exists {
// 			instanceTypes[instanceType] = true
// 		}

// 		// Fargate 작업인 경우, 인스턴스 유형으로 "fargate"를 추가
// 		// ProviderID가 "fargate"인 경우를 확인
// 		if node.Spec.ProviderID != "" && strings.Contains(node.Spec.ProviderID, "fargate") {
// 			instanceTypes["fargate"] = true
// 		}
// 	}

// 	// 인스턴스 유형들을 쉼표로 구분하여 출력
// 	var types []string
// 	for instanceType := range instanceTypes {
// 		types = append(types, instanceType)
// 	}
// 	fmt.Printf("Instance Types used in the cluster: %s\n", strings.Join(types, ", "))

// 	// 다양한 인스턴스를 사용하고 있는지 확인
// 	return len(instanceTypes) > 1
// }

package scalability

import (
	"context"
	"fmt"
	"strings"

	"eks-checklist/cmd/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckInstanceTypes checks if the cluster uses multiple instance types (e.g., for cost optimization, flexibility).
func CheckInstanceTypes(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "다양한 인스턴스 타입 사용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "클러스터에서 단일 인스턴스 타입만 사용 중입니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	instanceTypes := make(map[string]bool)

	for _, node := range nodes.Items {
		// 표준 인스턴스 타입 라벨 확인
		if instanceType, exists := node.Labels["beta.kubernetes.io/instance-type"]; exists {
			instanceTypes[instanceType] = true
			result.Resources = append(result.Resources, fmt.Sprintf("Node: %s | InstanceType: %s", node.Name, instanceType))
		}

		// Fargate 노드 포함 여부
		if node.Spec.ProviderID != "" && strings.Contains(node.Spec.ProviderID, "fargate") {
			instanceTypes["fargate"] = true
			result.Resources = append(result.Resources, fmt.Sprintf("Node: %s | InstanceType: fargate", node.Name))
		}
	}

	if len(instanceTypes) > 1 {
		result.Passed = true
	} else {
		result.Passed = false
	}

	return result
}

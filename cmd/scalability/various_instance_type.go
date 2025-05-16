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
		CheckName:  "[SCL-007] 다양한 인스턴스 타입 사용",
		Manual:     false,
		Passed:     true,
		FailureMsg: "클러스터에서 단일 인스턴스 타입만 사용 중입니다.",
		Runbook:    "https://fitcloud.github.io/eks-checklist/runbook/scalability/SCL-007",
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

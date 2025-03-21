package scalability

import (
	"context"
	"fmt"
	"log"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckNodeGroupUsage는 카펜터 전용 노드 그룹과 Fargate 사용 여부를 검사합니다.
func CheckNodeGroupUsage(client kubernetes.Interface) bool {
	// 클러스터의 노드를 가져옵니다.
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		log.Printf("노드 목록을 가져오는 중 오류 발생: %v", err)
		return false
	}

	// 카펜터 전용 노드 그룹 여부와 Fargate 사용 여부를 추적하는 변수
	var isCarpenterNodeGroup bool
	var isFargate bool

	// 노드 라벨을 통해 카펜터 전용 노드 그룹과 Fargate 사용 여부 확인
	for _, node := range nodes.Items {
		log.Printf("Checking node: %s, Labels: %v", node.Name, node.Labels)

		// 카펜터 전용 노드 그룹 여부 확인 (karpenter.sh/provisioner-name 라벨 확인)
		if _, found := node.Labels["karpenter.sh/provisioner-name"]; found {
			isCarpenterNodeGroup = true
		}

		// Fargate 사용 여부 확인 (eks.amazonaws.com/fargate-profile 라벨 확인)
		if profile, found := node.Labels["eks.amazonaws.com/fargate-profile"]; found && profile != "" {
			isFargate = true
		}
	}

	// 카펜터와 Fargate가 하나라도 true라면 성공
	if isCarpenterNodeGroup || isFargate {
		fmt.Println("PASS: Carpenter or Fargate node group is in use.")
		return true
	}

	// 둘 다 없으면 실패
	fmt.Println("FAIL: Neither Carpenter nor Fargate node group is in use.")
	return false
}

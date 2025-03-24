package stability

import (
	"context"
	"log"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Node의 label을 확인하여 "topology.kubernetes.io/zone"의 값을 기준으로 Multi-AZ 여부를 판단
func CheckNodeMultiAZ(client kubernetes.Interface) bool {
	// 모든 Node를 조회
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Printf("Failed to list nodes: %v", err)
		return false
	}

	// 중복을 제거할 Map
	uniqueZones := make(map[string]struct{})

	// 각 Node에서 "topology.kubernetes.io/zone" 라벨 확인
	for _, node := range nodes.Items {
		zone, exists := node.Labels["topology.kubernetes.io/zone"]
		if !exists {
			log.Printf("Node %s does not have a zone label", node.Name)
			continue // 라벨이 없으면 스킵
		}
		uniqueZones[zone] = struct{}{} // Map에 추가 (중복 제거)
	}

	// 서로 다른 AZ(zone)가 2개 이상이면 Multi-AZ 이므로 true 반환
	return len(uniqueZones) >= 2
}

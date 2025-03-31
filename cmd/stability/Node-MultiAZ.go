// 변경 전 코드

// package stability

// import (
// 	"context"
// 	"log"

// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // Node의 label을 확인하여 "topology.kubernetes.io/zone"의 값을 기준으로 Multi-AZ 여부를 판단
// func CheckNodeMultiAZ(client kubernetes.Interface) bool {
// 	// 모든 Node를 조회
// 	nodes, err := client.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
// 	if err != nil {
// 		log.Printf("Failed to list nodes: %v", err)
// 		return false
// 	}

// 	// 중복을 제거할 Map
// 	uniqueZones := make(map[string]struct{})

// 	// 각 Node에서 "topology.kubernetes.io/zone" 라벨 확인
// 	for _, node := range nodes.Items {
// 		zone, exists := node.Labels["topology.kubernetes.io/zone"]
// 		if !exists {
// 			log.Printf("Node %s does not have a zone label", node.Name)
// 			continue // 라벨이 없으면 스킵
// 		}
// 		uniqueZones[zone] = struct{}{} // Map에 추가 (중복 제거)
// 	}

// 	// 서로 다른 AZ(zone)가 2개 이상이면 Multi-AZ 이므로 true 반환
// 	return len(uniqueZones) >= 2
// }

// 변경 후 코드
package stability

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckNodeMultiAZ - 데이터 플레인 노드가 여러 가용영역(AZ)에 분산 배포되어 있는지 확인
func CheckNodeMultiAZ(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "다수의 가용 영역에 데이터 플레인 노드 배포",
		Manual:     false,
		Passed:     true,
		FailureMsg: "데이터 플레인 노드가 다수의 가용 영역(AZ)에 분산되어 있지 않습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	nodes, err := client.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	zoneMap := make(map[string][]string) // zone → node 목록

	for _, node := range nodes.Items {
		zone, exists := node.Labels["topology.kubernetes.io/zone"]
		if !exists {
			result.Resources = append(result.Resources, fmt.Sprintf("Node: %s (zone 라벨 없음)", node.Name))
			continue
		}
		zoneMap[zone] = append(zoneMap[zone], node.Name)
	}

	if len(zoneMap) >= 2 {
		result.Passed = true
		// result.SuccessMsg = fmt.Sprintf("데이터 플레인 노드가 %d개의 가용 영역에 분산되어 있습니다.", len(zoneMap))
		// for zone, nodes := range zoneMap {
		// 	result.Resources = append(result.Resources, fmt.Sprintf("Zone: %s | Nodes: %v", zone, nodes))
		// }
	} else {
		result.Passed = false
		result.FailureMsg = "데이터 플레인 노드가 다수의 가용 영역(AZ)에 분산되어 있지 않습니다."
		for zone, nodes := range zoneMap {
			result.Resources = append(result.Resources, fmt.Sprintf("Zone: %s | Nodes: %v", zone, nodes))
		}
	}

	return result
}

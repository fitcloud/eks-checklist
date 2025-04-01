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

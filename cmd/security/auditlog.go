// 기존 코드
// package security

// import (
// 	"fmt"

// 	"github.com/aws/aws-sdk-go-v2/service/eks/types"
// )

// // EksCluster 타입 정의
// type EksCluster struct {
// 	Cluster *types.Cluster
// }

// // audit 로그가 활성화되어 있는지 확인하는 함수
// func CheckAuditLoggingEnabled(eksCluster *EksCluster) bool {
// 	// Logging이 nil인지 확인
// 	if eksCluster.Cluster.Logging == nil {
// 		fmt.Println("No logging configuration found for this cluster.")
// 		return false
// 	}

// 	// ClusterLogging 배열을 순회하여 Enabled가 true인 경우에만 Types 배열을 확인
// 	for _, clusterLogging := range eksCluster.Cluster.Logging.ClusterLogging {
// 		if clusterLogging.Enabled != nil && *clusterLogging.Enabled {
// 			// Types 배열을 순회하여 audit 로그가 있는지 확인
// 			for _, logType := range clusterLogging.Types {
// 				// audit 로그가 있는 경우 true 반환
// 				if logType == "audit" {
// 					return true
// 				}
// 			}
// 		}
// 	}
// 	return false
// }

// 변경 후 코드
package security

import (
	"eks-checklist/cmd/common"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

// EksCluster 타입 정의
type EksCluster struct {
	Cluster *types.Cluster
}

// Audit 로그 활성화 여부를 체크하는 함수
func CheckAuditLoggingEnabled(eksCluster *EksCluster) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "Audit 로그 활성화",
		Manual:     false,
		Passed:     true,
		FailureMsg: "Audit 로그가 활성화되지 않았습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	if eksCluster.Cluster.Logging == nil {
		result.Passed = false
		result.FailureMsg = "클러스터에 로깅 설정이 없습니다."
		return result
	}

	auditLoggingEnabled := false
	for _, clusterLogging := range eksCluster.Cluster.Logging.ClusterLogging {
		if clusterLogging.Enabled != nil && *clusterLogging.Enabled {
			for _, logType := range clusterLogging.Types {
				if logType == "audit" {
					auditLoggingEnabled = true
					break
				}
			}
		}
	}

	if !auditLoggingEnabled {
		result.Resources = append(result.Resources, "Cluster: "+*eksCluster.Cluster.Name)
	}

	return result
}

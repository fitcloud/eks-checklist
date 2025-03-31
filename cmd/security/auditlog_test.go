package security_test

// import (
// 	"testing"

// 	"eks-checklist/cmd/security"
// 	"eks-checklist/cmd/testutils"

// 	"github.com/aws/aws-sdk-go-v2/service/eks/types"
// )

// func TestCheckAuditLoggingEnabled(t *testing.T) {
// 	testCases := testutils.LoadTestCases(t, "audit_logging.yaml")

// 	for _, tc := range testCases {
// 		name := tc["name"].(string)
// 		expectFailure := tc["expect_failure"].(bool)

// 		t.Run(name, func(t *testing.T) {
// 			var logSetup types.LogSetup

// 			// Enabled 필드 처리
// 			if enabled, ok := tc["enabled"]; ok && enabled != nil {
// 				b := enabled.(bool)
// 				logSetup.Enabled = &b
// 			}

// 			// Types 필드 처리
// 			typesList := []types.LogType{}
// 			for _, tRaw := range tc["types"].([]interface{}) {
// 				tStr := tRaw.(string)
// 				typesList = append(typesList, types.LogType(tStr))
// 			}
// 			logSetup.Types = typesList

// 			// 클러스터 생성
// 			cluster := &security.EksCluster{
// 				Cluster: &types.Cluster{
// 					Logging: &types.Logging{
// 						ClusterLogging: []types.LogSetup{logSetup},
// 					},
// 				},
// 			}

// 			// enabled가 아예 없을 경우 (nil)
// 			if tc["enabled"] == nil {
// 				cluster.Cluster.Logging = nil
// 			}

// 			result := security.CheckAuditLoggingEnabled(cluster)
// 			if result != !expectFailure {
// 				t.Errorf("Test '%s' failed: expected %v, got %v", name, !expectFailure, result)
// 			}
// 		})
// 	}
// }

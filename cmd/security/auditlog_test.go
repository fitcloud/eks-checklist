package security_test

import (
	"testing"

	"eks-checklist/cmd/security"
	"eks-checklist/cmd/testutils"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

func TestCheckAuditLoggingEnabled(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "audit_logging.yaml")

	for _, tc := range testCases {
		name := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)

		t.Run(name, func(t *testing.T) {
			var logSetup types.LogSetup

			// types 필드를 YAML에서 읽어와 []types.LogType으로 변환
			typesList := []types.LogType{}
			if tc["types"] != nil {
				for _, tRaw := range tc["types"].([]interface{}) {
					tStr := tRaw.(string)
					typesList = append(typesList, types.LogType(tStr))
				}
			}
			logSetup.Types = typesList

			// YAML의 enabled 필드 처리: 값이 nil이 아니라면 bool 값으로 할당
			if enabled, ok := tc["enabled"]; ok {
				if enabled != nil {
					boolEnabled := enabled.(bool)
					logSetup.Enabled = &boolEnabled
				}
			}

			// 클러스터 생성 (Name은 nil이 아니도록 설정)
			cluster := &security.EksCluster{
				Cluster: &types.Cluster{
					Name: strPtr(name),
					Logging: &types.Logging{
						ClusterLogging: []types.LogSetup{logSetup},
					},
				},
			}

			// enabled 필드가 nil이거나 types 배열이 빈 경우 로깅 설정 자체를 nil로 처리
			if tc["enabled"] == nil || len(typesList) == 0 {
				cluster.Cluster.Logging = nil
			}

			result := security.CheckAuditLoggingEnabled(cluster)
			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", name, expectPass, result.Passed)
			}
		})
	}
}

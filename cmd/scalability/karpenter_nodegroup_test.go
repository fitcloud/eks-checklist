package scalability

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"eks-checklist/cmd/testutils"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// TestCheckNodeGroupUsage는 YAML 파일("karpenter_nodegroup.yaml")에 있는 테스트 케이스를 읽어 실행하는 테스트 함수입니다.
func TestCheckNodeGroupUsage(t *testing.T) {
	// YAML 파일에서 테스트 케이스 로드
	testCases := testutils.LoadTestCases(t, "karpenter_nodegroup.yaml")

	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)

		// YAML에서 node_name과 node_labels 값을 읽습니다.
		nodeName, _ := tc["node_name"].(string)
		nodeLabels := []interface{}{}
		if nl, ok := tc["node_labels"]; ok && nl != nil {
			nodeLabels = nl.([]interface{})
		}

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// node_labels (예: "key:value")를 map으로 변환
			labelMap := make(map[string]string)
			for _, label := range nodeLabels {
				labelStr := fmt.Sprintf("%v", label)
				parts := splitLabel(labelStr)
				if len(parts) == 2 {
					labelMap[parts[0]] = parts[1]
				}
			}

			// 노드 생성
			_, err := client.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
				ObjectMeta: metaV1.ObjectMeta{
					Name:   nodeName,
					Labels: labelMap,
				},
			}, metaV1.CreateOptions{})
			if err != nil {
				t.Errorf("Failed to create node with label %v: %v", labelMap, err)
			}

			// 실제 검사 함수 호출 (CheckNodeGroupUsage는 common.CheckResult를 반환)
			result := CheckNodeGroupUsage(client)

			// 최종 결과 비교: result.Passed와 expectPass가 같아야 함.
			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result.Passed)
			}
		})
	}
}

// splitLabel 함수는 "key:value" 형식의 문자열을 ":"로 분리하여 슬라이스로 반환합니다.
func splitLabel(label string) []string {
	return strings.Split(label, ":")
}

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

// TestCheckNodeGroupUsage는 YAML 파일에 있는 테스트 케이스를 읽고 실행하는 테스트 함수
func TestCheckNodeGroupUsage(t *testing.T) {
	// YAML 파일에서 테스트 케이스 로드
	testCases := testutils.LoadTestCases(t, "karpenter_nodegroup.yaml")

	// 각 테스트 케이스 실행
	for _, test := range testCases {
		// 각 테스트 케이스에 대해 독립적으로 실행
		t.Run(test["name"].(string), func(t *testing.T) {
			// node_name과 node_labels를 test case에 해당하는 값으로 가져옵니다.
			nodeName := test["node_name"].(string)
			nodeLabels := test["node_labels"].([]interface{})
			expectFailure := test["expect_failure"].(bool)

			// Fake Kubernetes 클라이언트 생성
			client := fake.NewSimpleClientset()

			// 레이블을 map으로 변환
			labelMap := make(map[string]string)
			for _, label := range nodeLabels {
				// 레이블이 "key:value" 형식일 경우 "key"와 "value"로 분리
				labelParts := fmt.Sprintf("%v", label)
				parts := splitLabel(labelParts)
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

			// 실제 검사 함수 호출
			result := CheckNodeGroupUsage(client)

			// 검사 결과 비교
			// 결과가 성공인데 실패를 예상하면 실패
			if result && expectFailure {
				t.Errorf("Test %s failed: Expected failure, but it passed", test["name"].(string))
			}

			// 결과가 실패인데 성공을 예상하면 실패
			if !result && !expectFailure {
				t.Errorf("Test %s failed: Expected success, but it failed", test["name"].(string))
			}
		})
	}
}

// splitLabel 함수는 "key:value" 형식의 레이블을 분리하여 key, value를 반환합니다.
func splitLabel(label string) []string {
	// ":"로 구분하여 key와 value를 나눕니다.
	parts := strings.Split(label, ":")
	return parts
}

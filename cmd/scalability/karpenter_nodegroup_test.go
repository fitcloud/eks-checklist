package scalability

import (
	"context"
	utils "eks-checklist/cmd/utils"
	"testing"

	v1 "k8s.io/api/core/v1"                       // Node 타입을 가져오기 위해 사용
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1" // ObjectMeta를 가져오기 위해 사용
	"k8s.io/client-go/kubernetes/fake"
)

// TestCase는 각 테스트 케이스에 대한 구조체
type TestCase struct {
	Name          string   `yaml:"name"`
	NodeLabels    []string `yaml:"node_labels"`
	ExpectFailure bool     `yaml:"expect_failure"`
}

// TestCheckNodeGroupUsage는 YAML 파일에 있는 테스트 케이스를 읽고 실행하는 테스트 함수
func TestCheckNodeGroupUsage(t *testing.T) {
	// YAML 파일에서 테스트 케이스 로드
	testCases, err := utils.LoadTestCases("carpenter_node_group_test.yaml")
	if err != nil {
		t.Fatalf("Error loading test cases: %v", err)
	}

	// Fake Kubernetes 클라이언트 생성
	client := fake.NewSimpleClientset()

	// 각 테스트 케이스 실행
	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			// 가짜 노드 생성
			for _, label := range test.NodeLabels {
				client.CoreV1().Nodes().Create(context.TODO(), &v1.Node{
					ObjectMeta: metaV1.ObjectMeta{ // ObjectMeta는 metaV1 패키지에서 가져옴
						Name:   "node-1",
						Labels: map[string]string{label: label},
					},
				}, metaV1.CreateOptions{})
			}

			// 실제 검사 함수 호출
			isCarpenterNodeGroup, isFargate := CheckNodeGroupUsage(client)

			// 검사 결과 비교
			if isCarpenterNodeGroup && test.ExpectFailure {
				t.Errorf("Expected failure for Carpenter node group, but it was found")
			}

			if !isCarpenterNodeGroup && !test.ExpectFailure {
				t.Errorf("Expected success for Carpenter node group, but it was not found")
			}

			if isFargate && test.ExpectFailure {
				t.Errorf("Expected failure for Fargate profile, but it was found")
			}

			if !isFargate && !test.ExpectFailure {
				t.Errorf("Expected success for Fargate profile, but it was not found")
			}
		})
	}
}

package stability_test

import (
	"testing"

	"eks-checklist/cmd/common"
	"eks-checklist/cmd/stability"
	"eks-checklist/cmd/testutils"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
)

func TestCheckKarpenterNode(t *testing.T) {
	// YAML 파일 "karpenter_node.yaml"에서 테스트 케이스 로드
	testCases := testutils.LoadTestCases(t, "karpenter_node.yaml")

	// NodeClaim GVR 정의 (Karpenter v1beta1 기준)
	gvr := schema.GroupVersionResource{
		Group:    "karpenter.k8s.aws",
		Version:  "v1beta1",
		Resource: "nodeclaims",
	}
	// GVR에 대응되는 ListKind 등록
	listKinds := map[schema.GroupVersionResource]string{
		gvr: "NodeClaimList",
	}

	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)

		// YAML에서 karpenter_installed 값을 읽습니다.
		karpenterInstalled, ok := tc["karpenter_installed"].(bool)
		if !ok {
			karpenterInstalled = true // 기본값: 설치됨
		}

		// YAML에서 node_claim_present 값을 읽습니다.
		nodeClaimPresent, ok := tc["node_claim_present"].(bool)
		if !ok {
			nodeClaimPresent = false // 기본값: 없음
		}

		t.Run(testName, func(t *testing.T) {
			var objects []runtime.Object

			// nodeClaimPresent가 true이면 NodeClaim 객체 생성
			if nodeClaimPresent {
				nodeClaim := &unstructured.Unstructured{
					Object: map[string]interface{}{
						"apiVersion": "karpenter.k8s.aws/v1beta1",
						"kind":       "NodeClaim",
						"metadata": map[string]interface{}{
							"name": "test-nodeclaim",
						},
					},
				}
				objects = append(objects, nodeClaim)
			}

			// fake dynamic client 생성 (ListKind 등록 포함)
			client := fake.NewSimpleDynamicClientWithCustomListKinds(runtime.NewScheme(), listKinds, objects...)

			// Karpenter 설치 여부를 나타내는 CheckResult 전달
			karpenterCheck := common.CheckResult{Passed: karpenterInstalled}

			// 함수 호출: 두 인자(karpenterCheck, client) 전달
			result := stability.CheckKarpenterNode(karpenterCheck, client)

			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result.Passed)
			}
		})
	}
}

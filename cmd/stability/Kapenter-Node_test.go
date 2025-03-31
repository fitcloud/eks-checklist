package stability_test

// import (
// 	"testing"

// 	"eks-checklist/cmd/stability"
// 	"eks-checklist/cmd/testutils"

// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	"k8s.io/apimachinery/pkg/runtime"
// 	"k8s.io/apimachinery/pkg/runtime/schema"
// 	"k8s.io/client-go/dynamic/fake"
// )

// func TestCheckKarpenterNode(t *testing.T) {
// 	// YAML 파일로부터 테스트 케이스 로드
// 	testCases := testutils.LoadTestCases(t, "karpenter_node.yaml")

// 	// GVR 정의: karpenter nodeclaims
// 	gvr := schema.GroupVersionResource{
// 		Group:    "karpenter.k8s.aws",
// 		Version:  "v1beta1",
// 		Resource: "nodeclaims",
// 	}

// 	// GVR에 대응되는 ListKind 등록
// 	listKinds := map[schema.GroupVersionResource]string{
// 		gvr: "NodeClaimList",
// 	}

// 	for _, tc := range testCases {
// 		name := tc["name"].(string)
// 		expectFailure := tc["expect_failure"].(bool)

// 		t.Run(name, func(t *testing.T) {
// 			var objects []runtime.Object

// 			// NodeClaim 리소스를 생성할 케이스라면 unstructured 객체 생성
// 			if !expectFailure {
// 				nodeClaim := &unstructured.Unstructured{
// 					Object: map[string]interface{}{
// 						"apiVersion": "karpenter.k8s.aws/v1beta1",
// 						"kind":       "NodeClaim",
// 						"metadata": map[string]interface{}{
// 							"name": "test-nodeclaim",
// 						},
// 					},
// 				}
// 				objects = append(objects, nodeClaim)
// 			}

// 			// ✅ ListKind를 등록한 fake dynamic client 생성
// 			client := fake.NewSimpleDynamicClientWithCustomListKinds(
// 				runtime.NewScheme(),
// 				listKinds,
// 				objects...,
// 			)

// 			// 테스트 대상 함수 실행
// 			result := stability.CheckKarpenterNode(client)

// 			// 기대 결과와 비교
// 			if result != !expectFailure {
// 				t.Errorf("Test '%s' failed: expected %v, got %v", name, !expectFailure, result)
// 			}
// 		})
// 	}
// }

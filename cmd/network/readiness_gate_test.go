package network_test

// import (
// 	"context"
// 	"testing"

// 	"eks-checklist/cmd/network"
// 	"eks-checklist/cmd/testutils"

// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes/fake"
// )

// func TestCheckReadinessGateEnabled_YAML(t *testing.T) {
// 	// YAML 파일 "readiness_gate.yaml"에서 테스트 케이스 로드
// 	testCases := testutils.LoadTestCases(t, "readiness_gate.yaml")
// 	for _, tc := range testCases {
// 		testName := tc["name"].(string)
// 		expectPass := tc["expect_pass"].(bool)
// 		nsList := tc["namespaces"].([]interface{})

// 		t.Run(testName, func(t *testing.T) {
// 			client := fake.NewSimpleClientset()

// 			// YAML에 정의된 각 Namespace 객체 생성
// 			for _, ns := range nsList {
// 				nsDef := ns.(map[string]interface{})
// 				nsName := nsDef["name"].(string)

// 				// labels 필드는 map[string]interface{}로 로드되므로 map[string]string으로 변환
// 				labels := make(map[string]string)
// 				if rawLabels, ok := nsDef["labels"].(map[string]interface{}); ok {
// 					for key, value := range rawLabels {
// 						labels[key] = value.(string)
// 					}
// 				}

// 				nsObj := &corev1.Namespace{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:   nsName,
// 						Labels: labels,
// 					},
// 				}

// 				_, err := client.CoreV1().Namespaces().Create(context.TODO(), nsObj, metav1.CreateOptions{})
// 				if err != nil {
// 					t.Fatalf("Failed to create namespace %s: %v", nsName, err)
// 				}
// 			}

// 			// CheckReadinessGateEnabled 함수 실행 및 결과 검증
// 			result := network.CheckReadinessGateEnabled(client)
// 			if result.Passed != expectPass {
// 				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result.Passed)
// 			}
// 		})
// 	}
// }

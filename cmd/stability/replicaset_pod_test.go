package stability_test

// import (
// 	"context"
// 	"testing"

// 	"eks-checklist/cmd/stability"
// 	"eks-checklist/cmd/testutils"

// 	appsv1 "k8s.io/api/apps/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes/fake"
// )

// func TestPodReplicaSetCheck(t *testing.T) {
// 	// YAML 파일 "PodReplicaSetCheck.yaml"에서 테스트 케이스 로드
// 	testCases := testutils.LoadTestCases(t, "replicaset_pod.yaml")
// 	for _, tc := range testCases {
// 		testName, ok := tc["name"].(string)
// 		if !ok {
// 			t.Fatalf("Test case missing 'name' field")
// 		}

// 		// expected_failure가 true이면 실패가 예상됨 → 함수 반환값은 false여야 함.
// 		expectedFailure, ok := tc["expected_failure"].(bool)
// 		if !ok {
// 			t.Fatalf("Test case '%s' missing or invalid 'expected_failure' field", testName)
// 		}
// 		// 실제 기대값: !expected_failure
// 		expected := !expectedFailure

// 		replicaSetsRaw, ok := tc["replicasets"].([]interface{})
// 		if !ok {
// 			t.Fatalf("Test case '%s' missing 'replicasets' field", testName)
// 		}

// 		client := fake.NewSimpleClientset()

// 		// YAML에 정의된 각 ReplicaSet 객체 생성
// 		for _, rsRaw := range replicaSetsRaw {
// 			rsMap, ok := rsRaw.(map[string]interface{})
// 			if !ok {
// 				t.Fatalf("Test case '%s': replicaset is not a map", testName)
// 			}
// 			ns, ok := rsMap["namespace"].(string)
// 			if !ok {
// 				t.Fatalf("Test case '%s': replicaset missing 'namespace'", testName)
// 			}
// 			name, ok := rsMap["name"].(string)
// 			if !ok {
// 				t.Fatalf("Test case '%s': replicaset missing 'name'", testName)
// 			}
// 			// 먼저 필드 존재 여부 확인
// 			rawReplicas, exists := rsMap["replicas"]
// 			if !exists {
// 				t.Fatalf("Test case '%s': replicaset missing 'replicas'", testName)
// 			}
// 			var repValFloat float64
// 			switch v := rawReplicas.(type) {
// 			case float64:
// 				repValFloat = v
// 			case int:
// 				repValFloat = float64(v)
// 			default:
// 				t.Fatalf("Test case '%s': replicaset 'replicas' field has unexpected type: %T", testName, v)
// 			}
// 			replicaCount := int32(repValFloat)

// 			rsObj := &appsv1.ReplicaSet{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      name,
// 					Namespace: ns,
// 				},
// 				Spec: appsv1.ReplicaSetSpec{
// 					Replicas: &replicaCount,
// 				},
// 			}
// 			_, err := client.AppsV1().ReplicaSets(ns).Create(context.TODO(), rsObj, metav1.CreateOptions{})
// 			if err != nil {
// 				t.Fatalf("Test case '%s': failed to create ReplicaSet %s/%s: %v", testName, ns, name, err)
// 			}
// 		}

// 		// 함수 실행 및 반환값 비교
// 		result := stability.PodReplicaSetCheck(client)
// 		if result.Passed != expected {
// 			t.Errorf("Test '%s' failed: expected %v, got %v", testName, expected, result)
// 		}
// 	}
// }

package reliability_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/reliability"
	"eks-checklist/cmd/testutils"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestPodReplicaSetCheck(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "replicaset_pod.yaml")
	t.Logf("Loaded %d test cases", len(testCases)) // 로드된 케이스 수 확인

	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// YAML에 정의된 각 ReplicaSet 객체 생성
			replicaSetsRaw, ok := tc["replicasets"].([]interface{})
			if !ok {
				t.Fatalf("Test case '%s' missing 'replicasets' field", testName)
			}
			for _, rsRaw := range replicaSetsRaw {
				rsMap, ok := rsRaw.(map[string]interface{})
				if !ok {
					t.Fatalf("Test case '%s': replicaset is not a map", testName)
				}
				ns, ok := rsMap["namespace"].(string)
				if !ok {
					t.Fatalf("Test case '%s': replicaset missing 'namespace'", testName)
				}
				name, ok := rsMap["name"].(string)
				if !ok {
					t.Fatalf("Test case '%s': replicaset missing 'name'", testName)
				}
				rawReplicas, exists := rsMap["replicas"]
				if !exists {
					t.Fatalf("Test case '%s': replicaset missing 'replicas'", testName)
				}
				var repValFloat float64
				switch v := rawReplicas.(type) {
				case float64:
					repValFloat = v
				case int:
					repValFloat = float64(v)
				default:
					t.Fatalf("Test case '%s': replicaset 'replicas' field has unexpected type: %T", testName, v)
				}
				replicaCount := int32(repValFloat)

				rsObj := &appsv1.ReplicaSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: ns,
					},
					Spec: appsv1.ReplicaSetSpec{
						Replicas: &replicaCount,
					},
				}
				_, err := client.AppsV1().ReplicaSets(ns).Create(context.TODO(), rsObj, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Test case '%s': failed to create ReplicaSet %s/%s: %v", testName, ns, name, err)
				}
			}

			result := reliability.PodReplicaSetCheck(client)
			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result.Passed)
			}
		})
	}
}

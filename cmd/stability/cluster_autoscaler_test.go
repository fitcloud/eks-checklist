package stability_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/stability"
	"eks-checklist/cmd/testutils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1" // Container 등 사용
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckClusterAutoscalerEnabled(t *testing.T) {
	// YAML 파일 "CheckClusterAutoscalerEnabled.yaml"에서 테스트 케이스 로드
	testCases := testutils.LoadTestCases(t, "cluster_autoscaler.yaml")
	for _, tc := range testCases {
		testName, ok := tc["name"].(string)
		if !ok {
			t.Fatalf("Test case missing 'name' field")
		}
		// expected_failure: true이면 실패가 예상됨 → 함수 반환값은 false, 아니면 true.
		expectedFailure, ok := tc["expected_failure"].(bool)
		if !ok {
			t.Fatalf("Test case '%s' missing or invalid 'expected_failure' field", testName)
		}
		expected := !expectedFailure

		deploymentsRaw, ok := tc["deployments"].([]interface{})
		if !ok {
			t.Fatalf("Test case '%s' missing 'deployments' field", testName)
		}

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// YAML에 정의된 각 Deployment 객체 생성
			for _, depRaw := range deploymentsRaw {
				depMap, ok := depRaw.(map[string]interface{})
				if !ok {
					t.Fatalf("Test case '%s': deployment is not a map", testName)
				}
				ns, ok := depMap["namespace"].(string)
				if !ok {
					t.Fatalf("Test case '%s': deployment missing 'namespace'", testName)
				}
				name, ok := depMap["name"].(string)
				if !ok {
					t.Fatalf("Test case '%s': deployment missing 'name'", testName)
				}
				containersRaw, ok := depMap["containers"].([]interface{})
				if !ok {
					t.Fatalf("Test case '%s': deployment missing 'containers' field", testName)
				}
				var containers []corev1.Container
				for _, cRaw := range containersRaw {
					cMap, ok := cRaw.(map[string]interface{})
					if !ok {
						t.Fatalf("Test case '%s': container is not a map", testName)
					}
					cName, ok := cMap["name"].(string)
					if !ok {
						t.Fatalf("Test case '%s': container missing 'name'", testName)
					}
					image, ok := cMap["image"].(string)
					if !ok {
						t.Fatalf("Test case '%s': container missing 'image'", testName)
					}
					containers = append(containers, corev1.Container{
						Name:  cName,
						Image: image,
					})
				}

				depObj := &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: ns,
					},
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: containers,
							},
						},
					},
				}
				_, err := client.AppsV1().Deployments(ns).Create(context.TODO(), depObj, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Test case '%s': failed to create Deployment %s/%s: %v", testName, ns, name, err)
				}
			}

			// CheckClusterAutoscalerEnabled 함수 실행 및 반환값 검증
			result := stability.CheckClusterAutoscalerEnabled(client)
			if result != expected {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expected, result)
			}
		})
	}
}

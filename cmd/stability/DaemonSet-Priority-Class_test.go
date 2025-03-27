package stability_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/stability"
	"eks-checklist/cmd/testutils"

	corev1 "k8s.io/api/core/v1"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckDaemonSetPriorityClass(t *testing.T) {
	// YAML 파일 "CheckDaemonSetPriorityClass.yaml"에서 테스트 케이스 로드
	testCases := testutils.LoadTestCases(t, "daemonset_priority_class.yaml")
	for _, tc := range testCases {
		testName, ok := tc["name"].(string)
		if !ok {
			t.Fatalf("Test case missing 'name' field")
		}

		// expected_failure: true이면 실패가 예상됨 → 함수의 반환값은 false여야 함
		expectedFailureVal, ok := tc["expected_failure"]
		if !ok {
			t.Fatalf("Test case '%s' missing 'expected_failure' field", testName)
		}
		expectedFailure, ok := expectedFailureVal.(bool)
		if !ok {
			t.Fatalf("Test case '%s': expected_failure is not a bool", testName)
		}
		// 실제 기대값은 !expected_failure
		expected := !expectedFailure

		dsListRaw, ok := tc["daemonsets"].([]interface{})
		if !ok {
			t.Fatalf("Test case '%s' missing 'daemonsets' field", testName)
		}

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// YAML에 정의된 각 DaemonSet 객체 생성
			for _, dsRaw := range dsListRaw {
				dsMap, ok := dsRaw.(map[string]interface{})
				if !ok {
					t.Fatalf("Test case '%s': daemonset is not a map", testName)
				}
				// DaemonSet은 네임스페이스가 반드시 필요합니다.
				ns, ok := dsMap["namespace"].(string)
				if !ok {
					t.Fatalf("Test case '%s': daemonset missing 'namespace'", testName)
				}
				name, ok := dsMap["name"].(string)
				if !ok {
					t.Fatalf("Test case '%s': daemonset missing 'name'", testName)
				}
				priority, _ := dsMap["priorityClassName"].(string) // 없으면 빈 문자열

				dsObj := &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: ns,
					},
					Spec: appsv1.DaemonSetSpec{
						Template: corev1.PodTemplateSpec{ // metav1.PodTemplateSpec → corev1.PodTemplateSpec
							Spec: corev1.PodSpec{
								PriorityClassName: priority,
							},
						},
					},
				}
				_, err := client.AppsV1().DaemonSets(ns).Create(context.TODO(), dsObj, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Test case '%s': failed to create DaemonSet %s/%s: %v", testName, ns, name, err)
				}
			}

			// CheckDaemonSetPriorityClass 함수 실행 및 반환값 검증
			result := stability.CheckDaemonSetPriorityClass(client)
			if result != expected {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expected, result)
			}
		})
	}
}

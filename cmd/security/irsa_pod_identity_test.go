package security_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/security"
	"eks-checklist/cmd/testutils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckIRSAAndPodIdentity_YAML(t *testing.T) {
	// YAML 파일 "irsa_pod_identity.yaml"에서 테스트 케이스 로드
	testCases := testutils.LoadTestCases(t, "irsa_pod_identity.yaml")
	for _, tc := range testCases {
		testName := tc["name"].(string)
		// YAML의 "expected_pass" 필드를 그대로 사용
		expectedPass := tc["expected_pass"].(bool)

		serviceAccountsRaw, ok := tc["serviceaccounts"].([]interface{})
		if !ok {
			t.Fatalf("Test case '%s' is missing 'serviceaccounts' field", testName)
		}

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// YAML에 정의된 각 ServiceAccount 객체 생성 (kube-system 네임스페이스는 테스트 대상에서 제외됨)
			for _, saRaw := range serviceAccountsRaw {
				saMap, ok := saRaw.(map[string]interface{})
				if !ok {
					t.Fatalf("Test case '%s': serviceaccount is not a map", testName)
				}
				ns, ok := saMap["namespace"].(string)
				if !ok {
					t.Fatalf("Test case '%s': serviceaccount missing 'namespace'", testName)
				}
				name, ok := saMap["name"].(string)
				if !ok {
					t.Fatalf("Test case '%s': serviceaccount missing 'name'", testName)
				}
				annotations := make(map[string]string)
				if annRaw, exists := saMap["annotations"].(map[string]interface{}); exists {
					for key, val := range annRaw {
						sVal, ok := val.(string)
						if !ok {
							t.Fatalf("Test case '%s': annotation value for key %s is not a string", testName, key)
						}
						annotations[key] = sVal
					}
				}
				sa := &corev1.ServiceAccount{
					ObjectMeta: metav1.ObjectMeta{
						Name:        name,
						Namespace:   ns,
						Annotations: annotations,
					},
				}
				_, err := client.CoreV1().ServiceAccounts(ns).Create(context.TODO(), sa, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Test case '%s': failed to create serviceaccount %s/%s: %v", testName, ns, name, err)
				}
			}

			// CheckIRSAAndPodIdentity 함수 실행 및 반환값 검증
			result := security.CheckIRSAAndPodIdentity(client)
			if result.Passed != expectedPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectedPass, result)
			}
		})
	}
}

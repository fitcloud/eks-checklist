package network_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/network"
	"eks-checklist/cmd/testutils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckKubeProxyIPVSMode_YAML(t *testing.T) {
	// 생성한 YAML 파일 "kubeproxy_ipvs_mode.yaml"에서 테스트 케이스를 로드합니다.
	testCases := testutils.LoadTestCases(t, "kubeproxy_ipvs_mode.yaml")
	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)
		configValue := tc["config"].(string)

		t.Run(testName, func(t *testing.T) {
			// Fake Kubernetes 클라이언트 생성
			client := fake.NewSimpleClientset()

			// "kube-system" 네임스페이스에 kube-proxy ConfigMap 생성
			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kube-proxy-config",
					Namespace: "kube-system",
				},
				Data: map[string]string{
					"config": configValue,
				},
			}

			_, err := client.CoreV1().ConfigMaps("kube-system").Create(context.TODO(), configMap, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("Failed to create ConfigMap: %v", err)
			}

			// CheckKubeProxyIPVSMode 함수 실행 후 반환값 검증
			result := network.CheckKubeProxyIPVSMode(client)
			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result.Passed)
			}
		})
	}
}

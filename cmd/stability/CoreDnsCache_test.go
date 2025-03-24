package stability_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/stability"
	"eks-checklist/cmd/testutils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckCoreDNSCache(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "coredns_cache.yaml")

	for _, tc := range testCases {
		name := tc["name"].(string)
		expectFailure := tc["expect_failure"].(bool)

		t.Run(name, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			switch name {
			case "CoreDNS_with_cache":
				client.CoreV1().ConfigMaps("kube-system").Create(context.TODO(), &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: "coredns",
					},
					Data: map[string]string{
						"Corefile": ".:53 {\n    cache 30\n    forward . 8.8.8.8\n}\n",
					},
				}, metav1.CreateOptions{})

			case "CoreDNS_without_cache":
				client.CoreV1().ConfigMaps("kube-system").Create(context.TODO(), &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: "coredns",
					},
					Data: map[string]string{
						"Corefile": ".:53 {\n    forward . 8.8.8.8\n}\n",
					},
				}, metav1.CreateOptions{})

			case "Corefile_missing":
				client.CoreV1().ConfigMaps("kube-system").Create(context.TODO(), &corev1.ConfigMap{
					ObjectMeta: metav1.ObjectMeta{
						Name: "coredns",
					},
					Data: map[string]string{}, // Corefile key 없음
				}, metav1.CreateOptions{})

			case "ConfigMap_missing":
				// ConfigMap을 생성하지 않음
			}

			result := stability.CheckCoreDNSCache(client)

			if result != !expectFailure {
				t.Errorf("Test '%s' failed: expected %v, got %v", name, !expectFailure, result)
			}
		})
	}
}

package stability_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/stability"
	"eks-checklist/cmd/testutils"

	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckCoreDNSHpa(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "coredns_hpa.yaml")

	for _, tc := range testCases {
		name := tc["name"].(string)
		expectFailure := tc["expect_failure"].(bool)

		t.Run(name, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			if !expectFailure {
				// CoreDNS용 HPA 생성
				_, err := client.AutoscalingV1().HorizontalPodAutoscalers("kube-system").Create(context.TODO(), &autoscalingv1.HorizontalPodAutoscaler{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "coredns",
						Namespace: "kube-system",
					},
					Spec: autoscalingv1.HorizontalPodAutoscalerSpec{
						MaxReplicas: 5,
						ScaleTargetRef: autoscalingv1.CrossVersionObjectReference{
							Kind: "Deployment",
							Name: "coredns",
						},
					},
				}, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("failed to create HPA: %v", err)
				}
			}

			result := stability.CheckCoreDNSHpa(client)

			if result != !expectFailure {
				t.Errorf("Test '%s' failed: expected %v, got %v", name, !expectFailure, result)
			}
		})
	}
}

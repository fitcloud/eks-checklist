package scalability_test

import (
	"context"
	"strings"
	"testing"

	"eks-checklist/cmd/scalability"
	"eks-checklist/cmd/testutils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetKarpenter(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "karpenter_deploy.yaml")

	for _, tc := range testCases {
		name := tc["name"].(string)
		expectFailure := tc["expect_failure"].(bool)
		podImagesStr := tc["pod_images"].(string)

		t.Run(name, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// 주어진 이미지 목록으로 deployment 생성
			for _, img := range strings.Split(podImagesStr, ",") {
				_, err := client.AppsV1().Deployments("default").Create(context.TODO(), &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name: "test-deploy",
					},
					Spec: appsv1.DeploymentSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Name:  "container",
										Image: img,
									},
								},
							},
						},
					},
				}, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Failed to create test deployment: %v", err)
				}
			}

			result := scalability.GetKarpenter(client)

			if result != !expectFailure {
				t.Errorf("Test '%s' failed: expected %v, got %v", name, !expectFailure, result)
			}
		})
	}
}

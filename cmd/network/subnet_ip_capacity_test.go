package network_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/network"
	"eks-checklist/cmd/testutils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1" // Container 등 사용
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckAwsLoadBalancerController(t *testing.T) {
	// YAML 파일 "aws_loadbalancer_controller.yaml"에서 테스트 케이스를 로드합니다.
	testCases := testutils.LoadTestCases(t, "aws_loadbalancer_controller.yaml")
	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)
		deploymentsRaw := tc["deployments"].([]interface{})

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// YAML에 정의된 Deployment 객체들을 생성합니다.
			for _, d := range deploymentsRaw {
				dMap := d.(map[string]interface{})
				ns := dMap["namespace"].(string)
				name := dMap["name"].(string)

				containersRaw := dMap["containers"].([]interface{})
				var containers []corev1.Container
				for _, c := range containersRaw {
					cMap := c.(map[string]interface{})
					containers = append(containers, corev1.Container{
						Name:  cMap["name"].(string),
						Image: cMap["image"].(string),
					})
				}

				deployObj := &appsv1.Deployment{
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

				_, err := client.AppsV1().Deployments(ns).Create(context.TODO(), deployObj, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Failed to create Deployment %s/%s: %v", ns, name, err)
				}
			}

			// CheckAwsLoadBalancerController 함수 실행 후 반환값 검증
			result := network.CheckAwsLoadBalancerController(client)
			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result.Passed)
			}
		})
	}
}

package reliability_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/reliability"
	"eks-checklist/cmd/testutils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckNodeMultiAZ(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "node_multi_az.yaml")

	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			switch testName {
			case "Single_AZ":
				client.CoreV1().Nodes().Create(context.TODO(), &corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-1",
						Labels: map[string]string{
							"topology.kubernetes.io/zone": "us-east-1a",
						},
					},
				}, metav1.CreateOptions{})

				client.CoreV1().Nodes().Create(context.TODO(), &corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-2",
						Labels: map[string]string{
							"topology.kubernetes.io/zone": "us-east-1a",
						},
					},
				}, metav1.CreateOptions{})

			case "Multi_AZ":
				client.CoreV1().Nodes().Create(context.TODO(), &corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-1",
						Labels: map[string]string{
							"topology.kubernetes.io/zone": "us-east-1a",
						},
					},
				}, metav1.CreateOptions{})

				client.CoreV1().Nodes().Create(context.TODO(), &corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-2",
						Labels: map[string]string{
							"topology.kubernetes.io/zone": "us-east-1b",
						},
					},
				}, metav1.CreateOptions{})

			case "No_Zone_Label":
				client.CoreV1().Nodes().Create(context.TODO(), &corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "node-1",
						Labels: map[string]string{}, // 라벨 없음
					},
				}, metav1.CreateOptions{})
			}

			result := reliability.CheckNodeMultiAZ(client)

			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result.Passed)
			}
		})
	}
}

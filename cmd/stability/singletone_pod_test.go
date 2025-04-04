package stability_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/stability"
	"eks-checklist/cmd/testutils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestSingletonPodCheck(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "singleton_check.yaml")

	for _, tc := range testCases {
		name := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)

		t.Run(name, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			switch name {
			case "Deployment_with_1_Replica":
				replicas := int32(1)
				client.AppsV1().Deployments("default").Create(context.TODO(), &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name: "singleton-deploy",
					},
					Spec: appsv1.DeploymentSpec{
						Replicas: &replicas,
					},
				}, metav1.CreateOptions{})

			case "Standalone_Pod_Exists":
				client.CoreV1().Pods("default").Create(context.TODO(), &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: "standalone-pod",
					},
				}, metav1.CreateOptions{})

			case "StatefulSet_with_1_Replica":
				replicas := int32(1)
				client.AppsV1().StatefulSets("default").Create(context.TODO(), &appsv1.StatefulSet{
					ObjectMeta: metav1.ObjectMeta{
						Name: "singleton-statefulset",
					},
					Spec: appsv1.StatefulSetSpec{
						Replicas: &replicas,
					},
				}, metav1.CreateOptions{})

			case "Pod_with_NodeSelector":
				client.CoreV1().Pods("default").Create(context.TODO(), &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pod-with-nodeselector",
					},
					Spec: corev1.PodSpec{
						NodeSelector: map[string]string{
							"node-role.kubernetes.io/worker": "true",
						},
					},
				}, metav1.CreateOptions{})

			case "All_Good_No_Singletons":
				// intentionally empty — no problematic objects created
			}

			result := stability.SingletonPodCheck(client)

			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", name, expectPass, result.Passed)
			}
		})
	}
}

package scalability_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/scalability"
	"eks-checklist/cmd/testutils"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckInstanceTypes(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "instance_types.yaml")

	for _, tc := range testCases {
		name := tc["name"].(string)
		expectFailure := tc["expect_failure"].(bool)
		instanceTypes := tc["instance_types"].([]interface{})
		providerIDs := tc["provider_ids"].([]interface{})

		t.Run(name, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			for i := range instanceTypes {
				labels := map[string]string{}
				if instanceTypes[i].(string) != "" {
					labels["beta.kubernetes.io/instance-type"] = instanceTypes[i].(string)
				}

				node := &corev1.Node{
					ObjectMeta: v1.ObjectMeta{
						Name:   "node-" + name + "-" + string(rune(i)),
						Labels: labels,
					},
					Spec: corev1.NodeSpec{
						ProviderID: providerIDs[i].(string),
					},
				}
				client.CoreV1().Nodes().Create(context.TODO(), node, v1.CreateOptions{})
			}

			result := scalability.CheckInstanceTypes(client)

			if result != !expectFailure {
				t.Errorf("Test '%s' failed: expected %v, got %v", name, !expectFailure, result)
			}
		})
	}
}

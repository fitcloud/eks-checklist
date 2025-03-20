package stability_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

type StabilityTestCase struct {
	Name        string `yaml:"name"`
	Deployments []struct {
		Name     string `yaml:"name"`
		Replicas int    `yaml:"replicas"`
	} `yaml:"deployments"`
	Pods []struct {
		Name         string `yaml:"name"`
		OwnerRefs    int    `yaml:"owner_refs"`
		NodeSelector bool   `yaml:"node_selector"`
	} `yaml:"pods"`
	StatefulSets []struct {
		Name     string `yaml:"name"`
		Replicas int    `yaml:"replicas"`
	} `yaml:"statefulsets"`
	ExpectFailure bool `yaml:"expect_failure"`
}

func loadStabilityTestCases(t *testing.T) []StabilityTestCase {
	path, err := filepath.Abs(filepath.Join("..", "..", "testdata", "singleton_pod_check.yaml"))
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	var cases []StabilityTestCase
	if err := yaml.Unmarshal(data, &cases); err != nil {
		t.Fatalf("failed to parse YAML: %v", err)
	}

	return cases
}

func TestSingletonPodCheck(t *testing.T) {
	testCases := loadStabilityTestCases(t)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var objs []runtime.Object

			for _, d := range tc.Deployments {
				deployment := &appsv1.Deployment{
					ObjectMeta: v1.ObjectMeta{Name: d.Name},
					Spec:       appsv1.DeploymentSpec{Replicas: int32Ptr(int32(d.Replicas))},
				}
				objs = append(objs, deployment)
			}

			for _, s := range tc.StatefulSets {
				statefulSet := &appsv1.StatefulSet{
					ObjectMeta: v1.ObjectMeta{Name: s.Name},
					Spec:       appsv1.StatefulSetSpec{Replicas: int32Ptr(int32(s.Replicas))},
				}
				objs = append(objs, statefulSet)
			}

			for _, p := range tc.Pods {
				pod := &corev1.Pod{
					ObjectMeta: v1.ObjectMeta{Name: p.Name},
					Spec:       corev1.PodSpec{},
				}
				if p.NodeSelector {
					pod.Spec.NodeSelector = map[string]string{"key": "value"}
				}
				objs = append(objs, pod)
			}

			client := fake.NewSimpleClientset(objs...)
			result := SingletonPodCheck(client)
			assert.Equal(t, !tc.ExpectFailure, result, "Unexpected singleton pod check result")
		})
	}
}

func int32Ptr(i int32) *int32 { return &i }

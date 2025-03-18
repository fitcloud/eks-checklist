package general_test

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"

	general "eks-checklist/cmd/general"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakekube "k8s.io/client-go/kubernetes/fake"
)

type TestCase struct {
	Name          string   `yaml:"name"`
	PodImages     []string `yaml:"pod_images"`
	ExpectFailure bool     `yaml:"expect_failure"`
}

func loadTestCases(t *testing.T) []TestCase {
	path, err := filepath.Abs(filepath.Join("..", "..", "testdata", "check_image_tag.yaml"))
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	var cases []TestCase
	if err := yaml.Unmarshal(data, &cases); err != nil {
		t.Fatalf("failed to parse YAML: %v", err)
	}

	return cases
}

func TestCheckImageTag(t *testing.T) {
	cases := loadTestCases(t)

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// Fake client 생성
			client := fakekube.NewSimpleClientset()

			// Fake Pod 생성
			var podContainers []corev1.Container
			for _, img := range tc.PodImages {
				podContainers = append(podContainers, corev1.Container{Image: img})
			}

			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod", Namespace: "default"},
				Spec:       corev1.PodSpec{Containers: podContainers},
			}

			// Fake 클러스터에 Pod 추가
			_, err := client.CoreV1().Pods("default").Create(context.TODO(), pod, metav1.CreateOptions{})
			if err != nil {
				t.Fatalf("failed to create fake pod: %v", err)
			}

			// CheckImageTag 함수 실행 및 결과 검증
			result := general.CheckImageTag(client)
			if result != tc.ExpectFailure {
				t.Errorf("expected %v, got %v", tc.ExpectFailure, result)
			}
		})
	}
}

<<<<<<< Updated upstream
package general_test

import (
	"context"
	"testing"

	general "eks-checklist/cmd/general"
	"eks-checklist/cmd/testutils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakekube "k8s.io/client-go/kubernetes/fake"
)

func TestCheckImageTag(t *testing.T) {
	// testutils에서 LoadTestCases를 사용하여 YAML 파일 로드
	cases := testutils.LoadTestCases(t, "check_image_tag.yaml")

	// 각 테스트 케이스에 대해 반복
	for _, tc := range cases {
		// 'name' 필드를 가져와서 string 타입으로 처리
		name := tc["name"].(string)

		// 'expect_failure' 필드를 가져와서 bool 타입으로 처리
		expectFailure := tc["expect_failure"].(bool)

		// 'pod_images' 필드를 가져오고, []interface{}로 받은 후 []string으로 변환
		rawPodImages := tc["pod_images"].([]interface{})
		var podImages []string
		for _, img := range rawPodImages {
			podImages = append(podImages, img.(string))
		}

		// 테스트 실행
		t.Run(name, func(t *testing.T) {
			// Fake Kubernetes 클러스터 생성
			client := fakekube.NewSimpleClientset()

			// Fake Pod 생성
			var podContainers []corev1.Container
			for _, img := range podImages {
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
			if result != expectFailure {
				t.Errorf("expected %v, got %v", expectFailure, result)
			}
		})
	}
}
=======
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
>>>>>>> Stashed changes

package security_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/security"
	"eks-checklist/cmd/testutils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func boolPtr(b bool) *bool {
	return &b
}

func TestReadnonlyFilesystemCheck_YAML(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "readonly_filesystem.yaml")

	for _, tc := range testCases {
		name := tc["name"].(string)
		expectFailure := tc["expect_failure"].(bool)
		nodeOSList := tc["node_os"].([]interface{})
		pods := tc["pods"].([]interface{})

		// 윈도우 노드 테스트 케이스는 건너뜀
		skip := false
		for _, os := range nodeOSList {
			if os.(string) == "windows" {
				skip = true
				break
			}
		}
		if skip {
			t.Logf("Skipping test '%s' because it targets windows nodes", name)
			continue
		}

		t.Run(name, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// 노드 생성
			for i, osRaw := range nodeOSList {
				os := osRaw.(string)
				nodeName := "node-" + string(rune(i+1))
				client.CoreV1().Nodes().Create(context.TODO(), &corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name:   nodeName,
						Labels: map[string]string{"kubernetes.io/os": os},
					},
				}, metav1.CreateOptions{})
			}

			// Pod 생성
			for _, p := range pods {
				pdef := p.(map[string]interface{})
				containers := []corev1.Container{}

				for _, c := range pdef["containers"].([]interface{}) {
					cdef := c.(map[string]interface{})
					readOnly := boolPtr(false)
					if val, ok := cdef["readOnlyRootFilesystem"]; ok && val.(bool) {
						readOnly = boolPtr(true)
					}

					containers = append(containers, corev1.Container{
						Name: cdef["name"].(string),
						SecurityContext: &corev1.SecurityContext{
							ReadOnlyRootFilesystem: readOnly,
						},
					})
				}

				client.CoreV1().Pods(pdef["namespace"].(string)).Create(context.TODO(), &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pdef["name"].(string),
						Namespace: pdef["namespace"].(string),
					},
					Spec: corev1.PodSpec{
						NodeName:   pdef["node"].(string),
						Containers: containers,
					},
				}, metav1.CreateOptions{})
			}

			// 테스트 실행 및 결과 검증
			result := security.ReadnonlyFilesystemCheck(client)
			shouldPass := !expectFailure
			if result != shouldPass {
				t.Errorf("Test '%s' failed: expected pass = %v, got %v", name, shouldPass, result)
			}
		})
	}
}

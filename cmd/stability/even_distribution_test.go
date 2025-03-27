package stability_test

import (
	"bytes"
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"eks-checklist/cmd/stability"
	"eks-checklist/cmd/testutils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// captureOutput는 주어진 함수 실행 시 표준 출력을 캡쳐하여 문자열로 반환합니다.
func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestCheckPodDistributionAndAffinity(t *testing.T) {
	// YAML 파일 "CheckPodDistributionAndAffinity.yaml"에서 테스트 케이스 로드
	testCases := testutils.LoadTestCases(t, "even_distribution.yaml")
	for _, tc := range testCases {
		testName, ok := tc["name"].(string)
		if !ok {
			t.Fatalf("Test case missing 'name' field")
		}
		// expected_failure: true이면 실패가 예상됨 → 함수의 논리적 성공은 false, 그렇지 않으면 true.
		expectedFailure, ok := tc["expected_failure"].(bool)
		if !ok {
			t.Fatalf("Test case '%s' missing or invalid 'expected_failure' field", testName)
		}
		expected := !expectedFailure

		podsRaw, ok := tc["pods"].([]interface{})
		if !ok {
			t.Fatalf("Test case '%s' missing 'pods' field", testName)
		}

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// YAML에 정의된 각 Pod 객체 생성
			for _, pRaw := range podsRaw {
				pMap, ok := pRaw.(map[string]interface{})
				if !ok {
					t.Fatalf("Test case '%s': pod is not a map", testName)
				}
				ns, ok := pMap["namespace"].(string)
				if !ok {
					t.Fatalf("Test case '%s': pod missing 'namespace'", testName)
				}
				name, ok := pMap["name"].(string)
				if !ok {
					t.Fatalf("Test case '%s': pod missing 'name'", testName)
				}
				// affinity: if true → non-nil affinity; if false → nil
				var podAffinity *corev1.Affinity = nil
				if affinityVal, exists := pMap["affinity"]; exists {
					if aff, ok := affinityVal.(bool); ok && aff {
						// 실제 내용은 중요하지 않으므로 빈 구조체로 설정
						podAffinity = &corev1.Affinity{}
					}
				}
				// topologySpreadConstraints: YAML에서 배열로 지정, 각 항목에 "maxSkew" 필드만 사용
				var tscList []corev1.TopologySpreadConstraint
				if tscVal, exists := pMap["topologySpreadConstraints"]; exists {
					if arr, ok := tscVal.([]interface{}); ok {
						for _, item := range arr {
							itemMap, ok := item.(map[string]interface{})
							if !ok {
								continue
							}
							var maxSkew int32
							switch v := itemMap["maxSkew"].(type) {
							case float64:
								maxSkew = int32(v)
							case int:
								maxSkew = int32(v)
							default:
								maxSkew = 0
							}
							tscList = append(tscList, corev1.TopologySpreadConstraint{
								MaxSkew: maxSkew,
							})
						}
					}
				}

				podObj := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      name,
						Namespace: ns,
					},
					Spec: corev1.PodSpec{
						Affinity:                  podAffinity,
						TopologySpreadConstraints: tscList,
					},
				}

				_, err := client.CoreV1().Pods(ns).Create(context.TODO(), podObj, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Test case '%s': failed to create Pod %s/%s: %v", testName, ns, name, err)
				}
			}

			// CheckPodDistributionAndAffinity 함수 실행 및 표준 출력 캡쳐
			output := captureOutput(func() {
				stability.CheckPodDistributionAndAffinity(client)
			})

			// 기대값에 따라 "PASS" 또는 "FAIL" 문자열이 출력되었는지 확인
			if expected {
				if !strings.Contains(output, "PASS") {
					t.Errorf("Test '%s' failed: expected PASS in output, got:\n%s", testName, output)
				}
			} else {
				if !strings.Contains(output, "FAIL") {
					t.Errorf("Test '%s' failed: expected FAIL in output, got:\n%s", testName, output)
				}
			}
		})
	}
}

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

// int64Ptr 함수는 int64 값을 포인터로 반환합니다.
func int64Ptr(i int64) *int64 {
	return &i
}

// strPtr 함수는 문자열 값을 포인터로 반환합니다.
func strPtr(s string) *string {
	return &s
}

// TestCheckContainerExecutionUser_YAML 함수는 YAML 파일("container_user_check_test.yaml")에서
// 테스트 케이스를 로드하여 CheckContainerExecutionUser 함수가 올바르게 동작하는지 검증합니다.
func TestCheckContainerExecutionUser_YAML(t *testing.T) {
	// YAML 파일로부터 테스트 케이스를 로드합니다.
	testCases := testutils.LoadTestCases(t, "container_user_check_test.yaml")

	for _, tc := range testCases {
		name := tc["name"].(string)
		// expect_failure 값은 "루트 사용자(또는 Administrator) 감지 시" true여야 함.
		expectFailure := tc["expect_failure"].(bool)
		pods := tc["pods"].([]interface{})

		t.Run(name, func(t *testing.T) {
			// Fake Kubernetes 클라이언트 생성 (실제 클러스터 없이 테스트)
			client := fake.NewSimpleClientset()

			// YAML 파일에 정의된 각 Pod를 생성합니다.
			for _, p := range pods {
				pdef := p.(map[string]interface{})
				containers := []corev1.Container{}

				// Pod 내의 각 컨테이너에 대해 SecurityContext 설정
				for _, c := range pdef["containers"].([]interface{}) {
					cdef := c.(map[string]interface{})
					container := corev1.Container{
						Name:            cdef["name"].(string),
						SecurityContext: &corev1.SecurityContext{},
					}

					// runAsUser 설정 처리 (숫자 또는 null)
					if val, ok := cdef["runAsUser"]; ok {
						if val == nil {
							container.SecurityContext.RunAsUser = nil
						} else {
							switch v := val.(type) {
							case float64:
								container.SecurityContext.RunAsUser = int64Ptr(int64(v))
							case int:
								container.SecurityContext.RunAsUser = int64Ptr(int64(v))
							default:
								t.Fatalf("runAsUser의 타입이 예상과 다릅니다: %T", val)
							}
						}
					}

					// Windows 환경일 경우 windowsUser 설정 처리
					if winUser, ok := cdef["windowsUser"]; ok && winUser != nil {
						container.SecurityContext.WindowsOptions = &corev1.WindowsSecurityContextOptions{
							RunAsUserName: strPtr(winUser.(string)),
						}
					}

					containers = append(containers, container)
				}

				// Fake 클라이언트에 Pod 생성
				pod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      pdef["name"].(string),
						Namespace: pdef["namespace"].(string),
					},
					Spec: corev1.PodSpec{
						Containers: containers,
					},
				}
				_, err := client.CoreV1().Pods(pdef["namespace"].(string)).Create(context.TODO(), pod, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Pod 생성 실패: %v", err)
				}
			}

			// security 패키지의 CheckContainerExecutionUser 함수 실행
			result := security.CheckContainerExecutionUser(client)

			// 함수 결과와 YAML의 기대값(expect_failure)이 일치하는지 검증합니다.
			// (expect_failure가 true면 루트 사용자가 감지되어야 하므로 result도 true여야 함)
			if result != expectFailure {
				t.Errorf("테스트 '%s' 실패: 기대값 = %v, 실제값 = %v", name, expectFailure, result)
			}
		})
	}
}

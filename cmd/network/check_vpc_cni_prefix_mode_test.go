package network_test

// import (
// 	"context"
// 	"testing"

// 	"eks-checklist/cmd/network"
// 	"eks-checklist/cmd/testutils"

// 	appsv1 "k8s.io/api/apps/v1"
// 	corev1 "k8s.io/api/core/v1" // Pod, Container, EnvVar 등 사용
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes/fake"
// )

// func TestCheckVpcCniPrefixMode_YAML(t *testing.T) {
// 	// 생성한 YAML 파일 "vpc_cni_prefix_mode.yaml"로부터 테스트 케이스를 로드합니다.
// 	testCases := testutils.LoadTestCases(t, "vpc_cni_prefix_mode.yaml")
// 	for _, tc := range testCases {
// 		testName := tc["name"].(string)
// 		expectPass := tc["expect_pass"].(bool)
// 		daemonsets := tc["daemonsets"].([]interface{})

// 		t.Run(testName, func(t *testing.T) {
// 			client := fake.NewSimpleClientset()

// 			// YAML에 정의된 각 DaemonSet 객체 생성
// 			for _, ds := range daemonsets {
// 				dsDef := ds.(map[string]interface{})
// 				namespace := dsDef["namespace"].(string)
// 				name := dsDef["name"].(string)
// 				containersRaw := dsDef["containers"].([]interface{})
// 				var containers []corev1.Container

// 				// 각 컨테이너의 환경 변수(env) 정보 생성
// 				for _, c := range containersRaw {
// 					cDef := c.(map[string]interface{})
// 					contName := cDef["name"].(string)
// 					var envVars []corev1.EnvVar
// 					if envList, ok := cDef["env"].([]interface{}); ok {
// 						for _, e := range envList {
// 							eDef := e.(map[string]interface{})
// 							envVars = append(envVars, corev1.EnvVar{
// 								Name:  eDef["name"].(string),
// 								Value: eDef["value"].(string),
// 							})
// 						}
// 					}
// 					containers = append(containers, corev1.Container{
// 						Name: contName,
// 						Env:  envVars,
// 					})
// 				}

// 				// DaemonSet 객체 생성
// 				dsObj := &appsv1.DaemonSet{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name:      name,
// 						Namespace: namespace,
// 					},
// 					Spec: appsv1.DaemonSetSpec{
// 						Template: corev1.PodTemplateSpec{
// 							Spec: corev1.PodSpec{
// 								Containers: containers,
// 							},
// 						},
// 					},
// 				}

// 				_, err := client.AppsV1().DaemonSets(namespace).Create(context.TODO(), dsObj, metav1.CreateOptions{})
// 				if err != nil {
// 					t.Fatalf("Failed to create DaemonSet: %v", err)
// 				}
// 			}

// 			// CheckVpcCniPrefixMode 함수 실행 후 반환값 검증
// 			result := network.CheckVpcCniPrefixMode(client)
// 			if result != expectPass {
// 				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result)
// 			}
// 		})
// 	}
// }

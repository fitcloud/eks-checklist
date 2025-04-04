package stability_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/stability"
	"eks-checklist/cmd/testutils"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckProbe(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "probe_check.yaml")

	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			var pod corev1.Pod

			switch testName {
			case "All_Probes_Set":
				pod = corev1.Pod{
					ObjectMeta: v1.ObjectMeta{
						Name:      "all-probes-pod",
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:           "container",
								StartupProbe:   &corev1.Probe{},
								LivenessProbe:  &corev1.Probe{},
								ReadinessProbe: &corev1.Probe{},
							},
						},
					},
				}

			case "Missing_StartupProbe":
				pod = corev1.Pod{
					ObjectMeta: v1.ObjectMeta{
						Name:      "missing-startup",
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:           "container",
								LivenessProbe:  &corev1.Probe{},
								ReadinessProbe: &corev1.Probe{},
							},
						},
					},
				}

			case "Missing_LivenessProbe":
				pod = corev1.Pod{
					ObjectMeta: v1.ObjectMeta{
						Name:      "missing-liveness",
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:           "container",
								StartupProbe:   &corev1.Probe{},
								ReadinessProbe: &corev1.Probe{},
							},
						},
					},
				}

			case "Missing_ReadinessProbe":
				pod = corev1.Pod{
					ObjectMeta: v1.ObjectMeta{
						Name:      "missing-readiness",
						Namespace: "default",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:          "container",
								StartupProbe:  &corev1.Probe{},
								LivenessProbe: &corev1.Probe{},
							},
						},
					},
				}

			case "Pod_In_KubeSystem":
				pod = corev1.Pod{
					ObjectMeta: v1.ObjectMeta{
						Name:      "coredns",
						Namespace: "kube-system",
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name: "container",
							},
						},
					},
				}
			}

			// 파드 생성
			client.CoreV1().Pods(pod.Namespace).Create(context.TODO(), &pod, v1.CreateOptions{})

			// 함수 실행
			result := stability.CheckProbe(client)

			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectPass, result.Passed)
			}
		})
	}
}

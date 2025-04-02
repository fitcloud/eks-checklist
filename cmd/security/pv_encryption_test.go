package security_test

import (
	"context"
	"testing"

	"eks-checklist/cmd/security"
	"eks-checklist/cmd/testutils"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckPVEncryption(t *testing.T) {
	testCases := testutils.LoadTestCases(t, "pv_encryption.yaml")

	for _, tc := range testCases {
		testName := tc["name"].(string)
		expectPass := tc["expect_pass"].(bool)

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			switch testName {
			case "All_EBS_Encrypted":
				client.CoreV1().PersistentVolumes().Create(context.TODO(), &corev1.PersistentVolume{
					ObjectMeta: v1.ObjectMeta{Name: "pv-ebs-encrypted"},
					Spec: corev1.PersistentVolumeSpec{
						PersistentVolumeSource: corev1.PersistentVolumeSource{
							CSI: &corev1.CSIPersistentVolumeSource{
								Driver: "ebs.csi.aws.com",
								VolumeAttributes: map[string]string{
									"encrypted": "true",
								},
							},
						},
					},
				}, v1.CreateOptions{})

			case "Some_EBS_Unencrypted":
				client.CoreV1().PersistentVolumes().Create(context.TODO(), &corev1.PersistentVolume{
					ObjectMeta: v1.ObjectMeta{Name: "pv-ebs-encrypted"},
					Spec: corev1.PersistentVolumeSpec{
						PersistentVolumeSource: corev1.PersistentVolumeSource{
							CSI: &corev1.CSIPersistentVolumeSource{
								Driver: "ebs.csi.aws.com",
								VolumeAttributes: map[string]string{
									"encrypted": "true",
								},
							},
						},
					},
				}, v1.CreateOptions{})

				client.CoreV1().PersistentVolumes().Create(context.TODO(), &corev1.PersistentVolume{
					ObjectMeta: v1.ObjectMeta{Name: "pv-ebs-unencrypted"},
					Spec: corev1.PersistentVolumeSpec{
						PersistentVolumeSource: corev1.PersistentVolumeSource{
							CSI: &corev1.CSIPersistentVolumeSource{
								Driver: "ebs.csi.aws.com",
								VolumeAttributes: map[string]string{
									"encrypted": "false",
								},
							},
						},
					},
				}, v1.CreateOptions{})

			case "Non_EBS_PV":
				client.CoreV1().PersistentVolumes().Create(context.TODO(), &corev1.PersistentVolume{
					ObjectMeta: v1.ObjectMeta{Name: "pv-nfs"},
					Spec: corev1.PersistentVolumeSpec{
						PersistentVolumeSource: corev1.PersistentVolumeSource{
							HostPath: &corev1.HostPathVolumeSource{
								Path: "/mnt/data",
							},
						},
					},
				}, v1.CreateOptions{})
			}

			result := security.CheckPVEcryption(client)

			if result.Passed != expectPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, !expectPass, result.Passed)
			}
		})
	}
}

package security

import (
	"context"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckPVEcryption - Persistent Volume (PV)의 암호화 상태를 확인
func CheckPVEcryption(client kubernetes.Interface) bool {
	// 모든 PV를 조회
	pvs, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Println("Error retrieving persistent volumes:", err)
		return false
	}

	// 암호화 상태 확인용 변수
	allEncrypted := true

	// 각 PV에서 암호화 상태를 확인
	for _, pv := range pvs.Items {
		// PV의 스토리지 클래스가 EBS인 경우 암호화 여부 확인
		if pv.Spec.CSI != nil && pv.Spec.CSI.Driver == "ebs.csi.aws.com" {
			// EBS PV의 경우, 암호화 여부를 확인하려면 EBS CSI가 제공하는 암호화 여부를 확인해야 함
			// pv.Spec.CSI.VolumeAttributes["encrypted"] 필드를 확인하는 예시
			encrypted, exists := pv.Spec.CSI.VolumeAttributes["encrypted"]
			if exists && encrypted == "true" {
				log.Printf("PV %s is encrypted\n", pv.Name)
			} else {
				log.Printf("PV %s is NOT encrypted\n", pv.Name)
				allEncrypted = false
			}
		} else {
			// EBS가 아닌 경우, 암호화 여부를 확인할 수 없으므로 암호화되지 않았다고 가정
			// 실은 이건 문제가 있음 이것도 추후 고민
			log.Printf("PV %s is NOT encrypted\n", pv.Name)
			allEncrypted = false
		}
	}

	// 모든 PV가 암호화되어 있으므로 true 반환
	return allEncrypted
}

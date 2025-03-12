package stability

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CoreDNS의 HPA가 존재하는지 확인
func CheckCoreDNSHpa(client kubernetes.Interface) bool {
	// CoreDNS의 HPA는 kube-system 네임스페이스에 존재하므로 해당 네임스페이스에서 확인
	hpas, err := client.AutoscalingV1().HorizontalPodAutoscalers("kube-system").List(context.TODO(), v1.ListOptions{
		FieldSelector: "metadata.name=coredns", // CoreDNS라는 이름을 가진 HPA만 조회
	})

	if err != nil {
		panic(err.Error())
	}

	// CoreDNS의 HPA가 존재하는지 여부를 반환
	return len(hpas.Items) > 0
}

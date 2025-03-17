package security

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PrintAccessControl는 aws-auth ConfigMap을 확인하고 EKS 클러스터 정보를 출력합니다.
func PrintAccessControl(client kubernetes.Interface, clusterName string) bool {
	// 'aws-auth' ConfigMap을 'kube-system' 네임스페이스에서 가져옵니다.
	configMap, err := client.CoreV1().ConfigMaps("kube-system").Get(context.TODO(), "aws-auth", metav1.GetOptions{})
	if err != nil {
		log.Printf("aws-auth ConfigMap을 가져오는 중 오류 발생: %v", err)
		return false
	}

	// aws-auth ConfigMap의 내용을 출력합니다.
	fmt.Println("aws-auth ConfigMap:")
	for key, value := range configMap.Data {
		fmt.Printf("%s: %s\n", key, value)
	}

	// EKS 클러스터 정보 출력 (eksCluster를 사용)
	// 클러스터 이름은 eksCluster로 전달되며 추가적인 정보는 실제로 AWS SDK가 필요하지만,
	// 클러스터 이름만은 `eksCluster` 인수로 출력할 수 있습니다.
	fmt.Println("EKS 클러스터 정보:")
	fmt.Printf("클러스터 이름: %s\n", clusterName)
	// ARN과 엔드포인트 등의 정보는 AWS SDK를 통해 가져와야 하지만,
	// 예시로 하드코딩된 값을 사용할 수 있습니다.
	clusterArn := "arn:aws:eks:us-west-2:123456789012:cluster/" + clusterName // 클러스터 ARN 예시
	fmt.Printf("클러스터 ARN: %s\n", clusterArn)
	fmt.Printf("클러스터 엔드포인트: https://eks-cluster-endpoint\n") // 실제 엔드포인트는 AWS SDK를 통해 가져와야 함

	return true
}

package stability

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 클러스터에 Horizontal Pod Autoscaler가 설정정되어 있는지 확인
func CheckHpa(client kubernetes.Interface) {
	// 설정된 HPA가 있는 Deployment 목록
	deploymentsWithHPA := []string{}
	// 설정된 HPA가 없는 Deployment 목록
	deploymentsWithoutHPA := []string{}

	// 모든 Namespace에서 Deployment를 조회
	deployments, err := client.AppsV1().Deployments(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error fetching deployments: %v", err)
	}

	// 모든 Namespace에서 HPA를 조회
	hpas, err := client.AutoscalingV1().HorizontalPodAutoscalers(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error fetching HPAs: %v", err)
	}

	// 각 Deployment에 대해 Horizontal Pod Autoscaler 존재 여부 확인
	for _, deployment := range deployments.Items {
		hpaFound := false

		// 모든 HPA를 확인하여, 해당 Deployment에 대해 설정된 HPA가 있는지 확인
		for _, hpa := range hpas.Items {
			if hpa.Spec.ScaleTargetRef.Name == deployment.Name {
				hpaFound = true
				break
			}
		}

		// HPA가 있으면 HPA가 있는 목록에 추가
		if hpaFound {
			deploymentsWithHPA = append(deploymentsWithHPA, fmt.Sprintf("%s/%s", deployment.Namespace, deployment.Name))
		} else {
			// HPA가 없으면 HPA가 없는 목록에 추가
			deploymentsWithoutHPA = append(deploymentsWithoutHPA, fmt.Sprintf("%s/%s", deployment.Namespace, deployment.Name))
		}
	}

	// 결과 출력
	if len(deploymentsWithHPA) > 0 {
		fmt.Println("Deployments with Horizontal Pod Autoscaler:")
		for _, dep := range deploymentsWithHPA {
			fmt.Println(dep)
		}
	} else {
		fmt.Println("No deployments with Horizontal Pod Autoscaler found.")
	}

	if len(deploymentsWithoutHPA) > 0 {
		fmt.Println("\nDeployments without Horizontal Pod Autoscaler:")
		for _, dep := range deploymentsWithoutHPA {
			fmt.Println(dep)
		}
	} else {
		fmt.Println("No deployments without Horizontal Pod Autoscaler found.")
	}
}

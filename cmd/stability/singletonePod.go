package stability

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 클러스터에 Cluster Autoscaler가 설치되어 있는지 확인
func SingletonPodCheck(client kubernetes.Interface) bool {
	// kube-system 네임스페이스의 모든 Deployment 목록 가져오기

	result := true
	if checkDeploymentReplicas(client) {
		result = false
	}
	if checkStandalonePods(client) {
		result = false
	}
	if checkStatefulSetReplicas(client) {
		result = false
	}
	if checkNodeSelectorPods(client) {
		result = false
	}

	return result
}

// Deployment 중 replicas: 1 인 항목 체크
func checkDeploymentReplicas(client kubernetes.Interface) bool {
	deployments, err := client.AppsV1().Deployments("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	result := true

	for _, deployment := range deployments.Items {
		if deployment.Spec.Replicas != nil && *deployment.Spec.Replicas == 1 {
			fmt.Printf("[WARNING] Deployment %s in namespace %s has replicas: 1\n", deployment.Name, deployment.Namespace)
			result = false
		}
	}

	return result
}

// Standalone Pod 탐지 (Deployment, StatefulSet 등에 속하지 않은 Pod 찾기)
func checkStandalonePods(client kubernetes.Interface) bool {
	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	result := true

	for _, pod := range pods.Items {
		if pod.OwnerReferences == nil || len(pod.OwnerReferences) == 0 {
			fmt.Printf("[WARNING] Standalone Pod %s in namespace %s detected\n", pod.Name, pod.Namespace)
			result = false
		}
	}

	return result
}

// StatefulSet 중 replicas: 1 인 항목 체크
func checkStatefulSetReplicas(client kubernetes.Interface) bool {
	statefulSets, err := client.AppsV1().StatefulSets("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	result := true

	for _, statefulSet := range statefulSets.Items {
		if statefulSet.Spec.Replicas != nil && *statefulSet.Spec.Replicas == 1 {
			fmt.Printf("[WARNING] StatefulSet %s in namespace %s has replicas: 1\n", statefulSet.Name, statefulSet.Namespace)
			result = false
		}
	}

	return result
}

// 특정 노드에 강제 배치된 Pod 체크 (nodeSelector 사용 확인)
func checkNodeSelectorPods(client kubernetes.Interface) bool {
	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	result := true

	for _, pod := range pods.Items {
		if len(pod.Spec.NodeSelector) > 0 {
			fmt.Printf("[WARNING] Pod %s in namespace %s has nodeSelector set\n", pod.Name, pod.Namespace)
			result = false
		}
	}

	return result
}

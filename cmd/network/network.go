package network

import (
	"context"
	"log"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// EKS 클러스터가 사용하는 서브넷의 가용가능한 IP 주소가 10% 미만인 경우 SubnetId와 사용가능한 IP 개수를 반환하는 함수

// 모든 namespace에서 aws-node 데몬셋을 가져와서 containers.env.key 중 ENABLE_PREFIX_DELEGATION의 값이 true 인지 false인지 확인

// 모든 namespace에서 deployement를 가져와서 container.image에 aws-load-balancer-controller가 포함되어 있는지 확인
func CheckAwsLoadBalancerController(client kubernetes.Interface) bool {
	deploys, err := client.AppsV1().Deployments("").List(context.TODO(), v1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	for _, deploy := range deploys.Items {
		for _, container := range deploy.Spec.Template.Spec.Containers {
			if strings.Contains(container.Image, "aws-load-balancer-controller") {
				return true
			}
		}
	}

	return false
}

// ALB/NLB의 대상으로 Pod의 IP 주소를 사용하는지 확인

func CheckAwsLoadBalancerPodIp(client kubernetes.Interface) bool {
	// ingress class가 alb인 ingress를 모두 가져옴
	ingress, err := client.NetworkingV1().Ingresses("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	// annotation을 확인하여 target-type이 없거나 instacne일 경우 false 반환
	for _, ing := range ingress.Items {
		if ing.Annotations["alb.ingress.kubernetes.io/target-type"] == "" || ing.Annotations["alb.ingress.kubernetes.io/target-type"] == "instance" {
			// 디버깅 용
			log.Printf("ingress: %s", ing.Name)
			return false
		}
	}

	// service 객체 가져오기
	services, err := client.CoreV1().Services("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	// Service의 metadata.ownerReferences 필드에 Ingress 리소스가 포함되는지, 즉 aws-load-balancer-controller가 생성한 리소스인지 확인
	for _, svc := range services.Items {
		for _, ownerRef := range svc.ObjectMeta.OwnerReferences {
			if ownerRef.Kind == "Ingress" {
				// annotaition을 확인하여 target-type이 ip가 아닐 경우 false 반환
				if svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-nlb-target-type"] != "ip" {
					return false
				}
			}
		}
	}

	return true
}

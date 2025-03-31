// 변경 전 코드

// package network

// import (
// 	"context"
// 	"log"

// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // ALB/NLB의 대상으로 Pod의 IP 주소를 사용하는지 확인
// func CheckAwsLoadBalancerPodIp(client kubernetes.Interface) bool {
// 	// ingress class가 alb인 ingress를 모두 가져옴
// 	ingress, err := client.NetworkingV1().Ingresses("").List(context.TODO(), v1.ListOptions{})
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	// annotation을 확인하여 target-type이 없거나 instacne일 경우 false 반환
// 	for _, ing := range ingress.Items {
// 		if ing.Annotations["alb.ingress.kubernetes.io/target-type"] == "" || ing.Annotations["alb.ingress.kubernetes.io/target-type"] == "instance" {
// 			// 디버깅 용
// 			log.Printf("ingress: %s", ing.Name)
// 			return false
// 		}
// 	}

// 	// service 객체 가져오기
// 	services, err := client.CoreV1().Services("").List(context.TODO(), v1.ListOptions{})
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	// Service의 metadata.ownerReferences 필드에 Ingress 리소스가 포함되는지, 즉 aws-load-balancer-controller가 생성한 리소스인지 확인
// 	for _, svc := range services.Items {
// 		for _, ownerRef := range svc.ObjectMeta.OwnerReferences {
// 			if ownerRef.Kind == "Ingress" {
// 				// annotaition을 확인하여 target-type이 ip가 아닐 경우 false 반환
// 				if svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-nlb-target-type"] != "ip" {
// 					return false
// 				}
// 			}
// 		}
// 	}

// 	return true
// }

// 변경 후 코드
package network

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckAwsLoadBalancerPodIp checks whether ALB/NLB uses Pod IP as its target.
func CheckAwsLoadBalancerPodIp(controller_installed common.CheckResult, client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName: "ALB/NLB의 대상으로 Pod의 IP 사용",
		Manual:    false,
		Passed:    true,
		Runbook:   "https://your.runbook.url/latest-tag-image",
	}

	if !controller_installed.Passed {
		result.Passed = false
		result.FailureMsg = "AWS Load Balancer Controller가 설치되어 있지 않습니다"
		return result
	}

	hasFailure := false

	// 1. Ingress 체크 (ALB)
	ingresses, err := client.NetworkingV1().Ingresses("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	for _, ing := range ingresses.Items {
		targetType := ing.Annotations["alb.ingress.kubernetes.io/target-type"]
		if targetType == "" || targetType == "instance" {
			hasFailure = true
			result.Resources = append(result.Resources,
				fmt.Sprintf("Ingress: %s/%s | target-type: %s",
					ing.Namespace, ing.Name, targetType))
		}
	}

	// 2. Service 체크 (NLB)
	services, err := client.CoreV1().Services("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = fmt.Sprintf("Service 조회 실패: %v", err)
		return result
	}

	for _, svc := range services.Items {
		for _, ownerRef := range svc.OwnerReferences {
			if ownerRef.Kind == "Ingress" {
				target := svc.Annotations["service.beta.kubernetes.io/aws-load-balancer-nlb-target-type"]
				if target != "ip" {
					hasFailure = true
					result.Resources = append(result.Resources,
						fmt.Sprintf("Service: %s/%s | target-type: %s",
							svc.Namespace, svc.Name, target))
				}
			}
		}
	}

	if hasFailure {
		result.Passed = false
		result.FailureMsg = "일부 ALB/NLB 리소스가 Pod IP가 아닌 instance를 대상으로 사용하고 있습니다."
	} else {
		result.Passed = true
		// result.SuccessMsg = "모든 ALB/NLB가 Pod IP를 대상으로 사용하고 있습니다."
	}

	return result
}

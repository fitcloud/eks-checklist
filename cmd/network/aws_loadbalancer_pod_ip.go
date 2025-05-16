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
		CheckName: "[NET-006] ALB/NLB의 대상으로 Pod의 IP 사용",
		Manual:    false,
		Passed:    true,
		Runbook:   "https://fitcloud.github.io/eks-checklist/runbook/network/NET-006",
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

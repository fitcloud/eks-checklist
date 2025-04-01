package network

import (
	"context"
	"fmt"

	"eks-checklist/cmd/common"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// EndpointSlicesCheck checks whether all services use EndpointSlices instead of Endpoints.
func EndpointSlicesCheck(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName: "Endpoint 대신 EndpointSlices 사용",
		Manual:    false,
		Passed:    true,
		Runbook:   "https://your.runbook.url/latest-tag-image",
	}

	endpointSlices, err := client.DiscoveryV1().EndpointSlices("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	endpoints, err := client.CoreV1().Endpoints("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	// 서비스별 EndpointSlice 사용 여부 확인
	affectedServices := []string{}
	for _, ep := range endpoints.Items {
		serviceName := ep.Name
		namespace := ep.Namespace

		hasSlice := false
		for _, slice := range endpointSlices.Items {
			if slice.Namespace == namespace && slice.Labels["kubernetes.io/service-name"] == serviceName {
				hasSlice = true
				break
			}
		}

		if !hasSlice {
			affectedServices = append(affectedServices, fmt.Sprintf("Service: %s/%s", namespace, serviceName))
		}
	}

	if len(affectedServices) == 0 {
		result.Passed = true
		// result.SuccessMsg = "모든 서비스가 EndpointSlices를 사용하고 있습니다."
	} else {
		result.Passed = false
		result.FailureMsg = "일부 서비스가 아직 EndpointSlices 대신 Endpoints를 사용하고 있습니다."
		result.Resources = affectedServices
	}

	return result
}

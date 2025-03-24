package network

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	Red    = "\033[31m" // 빨간색
	Green  = "\033[32m" // 초록색
	Yellow = "\033[33m" // 노란색
	Reset  = "\033[0m"  // 기본 색상으로 리셋
)

// Endpoint Slices 사용 여부 검사
func EndpointSlicesCheck(client kubernetes.Interface) {
	// 모든 네임스페이스의 EndpointSlices 가져오기
	endpointSlices, err := client.DiscoveryV1().EndpointSlices("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 모든 네임스페이스의 Endpoints 가져오기
	endpoints, err := client.CoreV1().Endpoints("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 서비스별 Endpoint 사용 여부 확인
	var affectedServices []string

	for _, ep := range endpoints.Items {
		serviceName := ep.Name
		namespace := ep.Namespace

		// 같은 네임스페이스에 있는 EndpointSlices 확인
		hasSlice := false
		for _, slice := range endpointSlices.Items {
			if slice.Namespace == namespace && slice.Labels["kubernetes.io/service-name"] == serviceName {
				hasSlice = true
				break
			}
		}

		if !hasSlice {
			affectedServices = append(affectedServices, fmt.Sprintf("- %s/%s", namespace, serviceName))
		}
	}

	// 최종 결과 출력
	if len(affectedServices) == 0 {
		fmt.Println(Green + "✔ PASS:  All services in this cluster are using EndpointSlices" + Reset)
	} else {
		fmt.Println(Red + "✖ FAIL: Some services int this cluster are still using Endpoints" + Reset)
		fmt.Println("Affectred Resources:")
		for _, svc := range affectedServices {
			fmt.Println(svc)
		}
		fmt.Println("Runbook URL: https://fitcloud.github.io/eks-checklist/index.html")
	}
}

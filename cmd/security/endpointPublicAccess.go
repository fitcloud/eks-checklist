package security

import (
	"eks-checklist/cmd/common"
)

func CheckEndpointPublicAccess(eksCluster EksCluster) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "EKS 클러스터 API 엔드포인트 접근 제어(공인망, 사설망, IP 기반 제어)",
		Manual:     false,
		Passed:     true,
		SuccessMsg: "EKS 클러스터 API 엔드포인트 접근이 허용된 트래픽으로만 제한되어 있습니다.",
		FailureMsg: "EKS 클러스터 API 엔드포인트가 외부 공용 인터넷에서 접근 가능한 상태입니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	if eksCluster.Cluster.ResourcesVpcConfig.EndpointPublicAccess {
		result.Passed = false
		result.Resources = append(result.Resources, "클러스터 이름: "+*eksCluster.Cluster.Name)
	}

	return result
}

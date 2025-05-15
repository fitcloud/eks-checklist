package network

import (
	"eks-checklist/cmd/common"
)

// RegisterCheckers는 Network 카테고리의 모든 체크 항목을 등록합니다.
func RegisterCheckers(registry *common.CheckerRegistry, k8sClient interface{}, cfg interface{}, cluster string) {
	// VPC 서브넷에 충분한 IP 대역대 확보 - Automatic
	registry.RegisterFunc(
		"VPC 서브넷에 충분한 IP 대역대 확보",
		common.CategoryNetwork,
		func() common.CheckResult {
			return CheckVpcSubnetIpCapacity(cluster, cfg)
		},
	)

	// Pod에 부여할 IP 부족시 알림 설정 - Manual
	registry.RegisterFunc(
		"Pod에 부여할 IP 부족시 알림 설정",
		common.CategoryNetwork,
		func() common.CheckResult {
			return CheckPodIpAlert()
		},
	)

	// VPC CNI의 Prefix 모드 사용 - Automatic
	registry.RegisterFunc(
		"VPC CNI의 Prefix 모드 사용",
		common.CategoryNetwork,
		func() common.CheckResult {
			return CheckVpcCniPrefixMode(k8sClient)
		},
	)

	// 사용 사례에 맞는 로드밸런서 사용(ALB or NLB) - Manual
	registry.RegisterFunc(
		"사용 사례에 맞는 로드밸런서 사용(ALB or NLB)",
		common.CategoryNetwork,
		func() common.CheckResult {
			return CheckLoadBalancerUsage()
		},
	)

	// AWS Load Balancer Controller 사용 - Automatic
	registry.RegisterFunc(
		"AWS Load Balancer Controller 사용",
		common.CategoryNetwork,
		func() common.CheckResult {
			return CheckALBController(k8sClient)
		},
	)

	// ALB/NLB의 대상으로 Pod의 IP 사용 - Automatic
	registry.RegisterFunc(
		"ALB/NLB의 대상으로 Pod의 IP 사용",
		common.CategoryNetwork,
		func() common.CheckResult {
			return CheckPodIpTarget(k8sClient)
		},
	)

	// Pod Readiness Gate 적용 - Automatic
	registry.RegisterFunc(
		"Pod Readiness Gate 적용",
		common.CategoryNetwork,
		func() common.CheckResult {
			return CheckReadinessGate(k8sClient)
		},
	)

	// kube-proxy에 IPVS 모드 적용 - Automatic
	registry.RegisterFunc(
		"kube-proxy에 IPVS 모드 적용",
		common.CategoryNetwork,
		func() common.CheckResult {
			return CheckKubeProxyIpvsMode(k8sClient)
		},
	)

	// Endpoint 대신 EndpointSlices 사용 - Automatic
	registry.RegisterFunc(
		"Endpoint 대신 EndpointSlices 사용",
		common.CategoryNetwork,
		func() common.CheckResult {
			return CheckEndpointSlices(k8sClient)
		},
	)
} 
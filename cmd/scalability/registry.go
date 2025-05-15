package scalability

import (
	"eks-checklist/cmd/common"
)

// RegisterCheckers는 Scalability 카테고리의 모든 체크 항목을 등록합니다.
func RegisterCheckers(registry *common.CheckerRegistry, k8sClient interface{}, cfg interface{}, cluster string) {
	// Karpenter 사용 - Automatic
	registry.RegisterFunc(
		"Karpenter 사용",
		common.CategoryScalability,
		func() common.CheckResult {
			return GetKarpenter(k8sClient)
		},
	)

	// Karpenter 전용 노드 그룹 혹은 Fargate 사용 - Automatic
	registry.RegisterFunc(
		"Karpenter 전용 노드 그룹 혹은 Fargate 사용",
		common.CategoryScalability,
		func() common.CheckResult {
			return CheckNodeGroupUsage(k8sClient)
		},
	)

	// Spot 노드 사용시 Spot 중지 핸들러 적용 - Automatic
	registry.RegisterFunc(
		"Spot 노드 사용시 Spot 중지 핸들러 적용",
		common.CategoryScalability,
		func() common.CheckResult {
			return CheckSpotNodeTerminationHandler(k8sClient)
		},
	)

	// 중요 Pod에 노드 삭제 방지용 Label 부여 - Manual
	registry.RegisterFunc(
		"중요 Pod에 노드 삭제 방지용 Label 부여",
		common.CategoryScalability,
		func() common.CheckResult {
			return CheckImportantPodProtection(k8sClient, cfg, cluster)
		},
	)

	// Application에 Graceful shutdown 적용 - Manual
	registry.RegisterFunc(
		"Application에 Graceful shutdown 적용",
		common.CategoryScalability,
		func() common.CheckResult {
			return CheckGracefulShutdown(k8sClient, cfg, cluster)
		},
	)

	// 노드 확장/축소 정책 적용 - Manual
	registry.RegisterFunc(
		"노드 확장/축소 정책 적용",
		common.CategoryScalability,
		func() common.CheckResult {
			return CheckNodeScalingPolicy()
		},
	)

	// 다양한 인스턴스 타입 사용 - Automatic
	registry.RegisterFunc(
		"다양한 인스턴스 타입 사용",
		common.CategoryScalability,
		func() common.CheckResult {
			return CheckInstanceTypes(k8sClient)
		},
	)
} 
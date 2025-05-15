package cost

import (
	"eks-checklist/cmd/common"
)

// RegisterCheckers는 Cost 카테고리의 모든 체크 항목을 등록합니다.
func RegisterCheckers(registry *common.CheckerRegistry, k8sClient interface{}, cfg interface{}, cluster string) {
	// kubecost 설치 - Automatic
	registry.RegisterFunc(
		"kubecost 설치",
		common.CategoryCost,
		func() common.CheckResult {
			return CheckKubecostInstalled(k8sClient)
		},
	)
} 
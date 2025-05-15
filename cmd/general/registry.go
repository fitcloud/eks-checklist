package general

import (
	"eks-checklist/cmd/common"
)

// RegisterCheckers는 General 카테고리의 모든 체크 항목을 등록합니다.
func RegisterCheckers(registry *common.CheckerRegistry, k8sClient interface{}) {
	// 코드형 인프라 체크
	registry.RegisterFunc(
		"코드형 인프라 (EKS 클러스터, 애플리케이션 배포)",
		common.CategoryGeneral,
		func() common.CheckResult {
			return CheckIAC()
		},
	)

	// GitOps 적용 체크
	registry.RegisterFunc(
		"GitOps 적용",
		common.CategoryGeneral,
		func() common.CheckResult {
			return CheckGitOps()
		},
	)

	// 컨테이너 이미지 태그에 latest 미사용 체크
	registry.RegisterFunc(
		"컨테이너 이미지 태그에 latest 미사용",
		common.CategoryGeneral,
		func() common.CheckResult {
			return CheckImageTag(k8sClient)
		},
	)
} 
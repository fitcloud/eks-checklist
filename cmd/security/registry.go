package security

import (
	"eks-checklist/cmd/common"
)

// RegisterCheckers는 Security 카테고리의 모든 체크 항목을 등록합니다.
func RegisterCheckers(registry *common.CheckerRegistry, k8sClient interface{}, cfg interface{}, cluster string, eksCluster interface{}) {
	// EKS 클러스터 API 엔드포인트 접근 제어(공인망, 사설망, IP 기반 제어) - Automatic
	registry.RegisterFunc(
		"EKS 클러스터 API 엔드포인트 접근 제어",
		common.CategorySecurity,
		func() common.CheckResult {
			return CheckEndpointPublicAccess(EksCluster(eksCluster))
		},
	)

	// 클러스터 접근 제어(Access entries, aws-auth 컨피그맵) - Automatic/Manual
	registry.RegisterFunc(
		"클러스터 접근 제어(Access entries, aws-auth 컨피그맵)",
		common.CategorySecurity,
		func() common.CheckResult {
			return CheckAccessControl(k8sClient, cfg, cluster)
		},
	)

	// IRSA 또는 Pod Identity 기반 권한 부여 - Automatic
	registry.RegisterFunc(
		"IRSA 또는 Pod Identity 기반 권한 부여",
		common.CategorySecurity,
		func() common.CheckResult {
			return CheckIRSAAndPodIdentity(k8sClient)
		},
	)

	// 데이터 플레인 노드에 필수로 필요한 IAM 권한만 부여 - Automatic
	registry.RegisterFunc(
		"데이터 플레인 노드에 필수로 필요한 IAM 권한만 부여",
		common.CategorySecurity,
		func() common.CheckResult {
			return CheckNodeIAMRoles(k8sClient)
		},
	)

	// 루트 유저가 아닌 유저로 컨테이너 실행 - Automatic
	registry.RegisterFunc(
		"루트 유저가 아닌 유저로 컨테이너 실행",
		common.CategorySecurity,
		func() common.CheckResult {
			return CheckContainerExecutionUser(k8sClient)
		},
	)

	// 멀티 태넌시 적용 유무 - Manual
	registry.RegisterFunc(
		"멀티 태넌시 적용 유무",
		common.CategorySecurity,
		func() common.CheckResult {
			return CheckMultitenancy(k8sClient, cfg, cluster)
		},
	)

	// Audit 로그 활성화 - Automatic
	registry.RegisterFunc(
		"Audit 로그 활성화",
		common.CategorySecurity,
		func() common.CheckResult {
			return CheckAuditLoggingEnabled(&EksCluster{Cluster: eksCluster.Cluster})
		},
	)

	// 비정상 접근에 대한 알림 설정 - Manual
	registry.RegisterFunc(
		"비정상 접근에 대한 알림 설정",
		common.CategorySecurity,
		func() common.CheckResult {
			return CheckAccessAlarm()
		},
	)

	// Pod-to-Pod 접근 제어 - Automatic/Manual
	registry.RegisterFunc(
		"Pod-to-Pod 접근 제어",
		common.CategorySecurity,
		func() common.CheckResult {
			return CheckPodToPodNetworkPolicy(k8sClient, cluster)
		},
	)

	// PV 암호화 - Automatic
	registry.RegisterFunc(
		"PV 암호화",
		common.CategorySecurity,
		func() common.CheckResult {
			return CheckPVEcryption(k8sClient)
		},
	)

	// Secret 객체 암호화 - Automatic
	registry.RegisterFunc(
		"Secret 객체 암호화",
		common.CategorySecurity,
		func() common.CheckResult {
			return CheckSecretEncryption(k8sClient)
		},
	)

	// 데이터 플레인 사설망 - Automatic
	registry.RegisterFunc(
		"데이터 플레인 사설망",
		common.CategorySecurity,
		func() common.CheckResult {
			return DataplanePrivateCheck(EksCluster(eksCluster), cfg)
		},
	)

	// 컨테이너 이미지 정적 분석 - Manual
	registry.RegisterFunc(
		"컨테이너 이미지 정적 분석",
		common.CategorySecurity,
		func() common.CheckResult {
			return CheckImageStaticAnalysis(k8sClient, cfg, cluster)
		},
	)

	// 읽기 전용 파일시스템 사용 - Automatic
	registry.RegisterFunc(
		"읽기 전용 파일시스템 사용",
		common.CategorySecurity,
		func() common.CheckResult {
			return ReadnonlyFilesystemCheck(k8sClient)
		},
	)
} 
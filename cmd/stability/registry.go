package stability

import (
	"eks-checklist/cmd/common"
	"eks-checklist/cmd/scalability"
)

// RegisterCheckers는 Stability 카테고리의 모든 체크 항목을 등록합니다.
func RegisterCheckers(registry *common.CheckerRegistry, k8sClient interface{}, cfg interface{}, cluster string, dynamicClient interface{}) {
	// 싱글톤 Pod 미사용 - Automatic
	registry.RegisterFunc(
		"싱글톤 Pod 미사용",
		common.CategoryStability,
		func() common.CheckResult {
			return SingletonPodCheck(k8sClient)
		},
	)

	// 2개 이상의 Pod 복제본 사용 - Automatic
	registry.RegisterFunc(
		"2개 이상의 Pod 복제본 사용",
		common.CategoryStability,
		func() common.CheckResult {
			return PodReplicaSetCheck(k8sClient)
		},
	)

	// 동일한 역할을 하는 Pod를 다수의 노드에 분산 배포 - Automatic
	registry.RegisterFunc(
		"동일한 역할을 하는 Pod를 다수의 노드에 분산 배포",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckPodDistributionAndAffinity(k8sClient)
		},
	)

	// HPA 적용 - Automatic
	registry.RegisterFunc(
		"HPA 적용",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckHpa(k8sClient)
		},
	)

	// Probe(Startup, Readiness, Liveness) 적용 - Automatic
	registry.RegisterFunc(
		"Probe(Startup, Readiness, Liveness) 적용",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckProbe(k8sClient)
		},
	)

	// 중요 워크로드에 대한 PDB(Pod Distruption Budget) 적용 - Automatic/Manual
	registry.RegisterFunc(
		"중요 워크로드에 대한 PDB(Pod Distruption Budget) 적용",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckPDB()
		},
	)

	// 애플리케이션에 적절한 CPU/RAM 할당 - Automatic/Manual
	registry.RegisterFunc(
		"애플리케이션에 적절한 CPU/RAM 할당",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckResourceAllocation(k8sClient, cfg, cluster)
		},
	)

	// 애플리케이션 중요도에 따른 QoS 적용 - Automatic/Manual
	registry.RegisterFunc(
		"애플리케이션 중요도에 따른 QoS 적용",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckQoSClass(k8sClient, cfg, cluster)
		},
	)

	// 인프라 및 애플리케이션 모니터링 스택 적용 - Manual
	registry.RegisterFunc(
		"인프라 및 애플리케이션 모니터링 스택 적용",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckNodeScalingPolicy()
		},
	)

	// 반영구 저장소에 애플리케이션 로그 저장 - Manual
	registry.RegisterFunc(
		"반영구 저장소에 애플리케이션 로그 저장",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckApplicationLogs()
		},
	)

	// 오토스케일링 그룹 기반 관리형 노드 그룹 생성 - Automatic
	registry.RegisterFunc(
		"오토스케일링 그룹 기반 관리형 노드 그룹 생성",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckAutoScaledManagedNodeGroup(k8sClient, cluster)
		},
	)

	// Cluster Autoscaler 적용 - Automatic
	registry.RegisterFunc(
		"Cluster Autoscaler 적용",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckClusterAutoscalerEnabled(k8sClient)
		},
	)

	// Karpenter 기반 노드 생성 - Automatic
	registry.RegisterFunc(
		"Karpenter 기반 노드 생성",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckKarpenterNode(scalability.GetKarpenter(k8sClient), dynamicClient)
		},
	)

	// 다수의 가용 영역에 데이터 플레인 노드 배포 - Automatic
	registry.RegisterFunc(
		"다수의 가용 영역에 데이터 플레인 노드 배포",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckNodeMultiAZ(k8sClient)
		},
	)

	// PV 사용시 volume affinity 위반 사항 체크 - Automatic/Manual
	registry.RegisterFunc(
		"PV 사용시 volume affinity 위반 사항 체크",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckVolumeAffinity()
		},
	)

	// CoreDNS에 HPA 적용 - Automatic
	registry.RegisterFunc(
		"CoreDNS에 HPA 적용",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckCoreDNSHPA(k8sClient)
		},
	)

	// DNS 캐시 적용 - Automatic/Manual
	registry.RegisterFunc(
		"DNS 캐시 적용",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckDNSCaching(k8sClient)
		},
	)

	// Karpenter 사용시 DaemonSet에 Priority Class 부여 - Automatic
	registry.RegisterFunc(
		"Karpenter 사용시 DaemonSet에 Priority Class 부여",
		common.CategoryStability,
		func() common.CheckResult {
			return CheckDaemonSetPriorityClass(k8sClient)
		},
	)
} 
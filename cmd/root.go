package cmd

import (
	"eks-checklist/cmd/common"
	"eks-checklist/cmd/cost"
	"eks-checklist/cmd/general"
	"eks-checklist/cmd/network"
	"eks-checklist/cmd/scalability"
	"eks-checklist/cmd/security"
	"eks-checklist/cmd/stability"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	kubeconfigPath    string
	kubeconfigContext string
	awsProfile        string
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "eks-checklist",
	Long:  "eks-checklist",
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig := getKubeconfig(kubeconfigPath, awsProfile)
		cluster := getEksClusterName(kubeconfig)

		fmt.Printf("Running checks on %s\n", cluster)

		eksCluster := Describe(cluster)
		k8sClient := createK8sClient(kubeconfig)
		cfg := GetAWSConfig()
		dynamicClient, err := CreateDynamicClient(&kubeconfig)
		if err != nil {
			fmt.Println("Error creating dynamic client:", err)
			os.Exit(1)
		}

		// General 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[General Check]===============\n")

		// 코드형 인프라 (EKS 클러스터, 애플리케이션 배포)
		common.PrintResult(general.CheckIAC())

		// GitOps 적용
		common.PrintResult(general.CheckGitOps())

		// 컨테이너 이미지 태그에 latest 미사용
		common.PrintResult(general.CheckImageTag(k8sClient))

		// Security 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Security Check]===============\n")

		// EKS 클러스터 API 엔드포인트 접근 제어(공인망, 사설망, IP 기반 제어) - Automatic
		common.PrintResult(security.CheckEndpointPublicAccess(security.EksCluster(eksCluster)))

		// 클러스터 접근 제어(Access entries, aws-auth 컨피그맵) - Automatic/Manual
		common.PrintResult(security.CheckAccessControl(k8sClient, cfg, cluster))

		// IRSA 또는 Pod Identity 기반 권한 부여 - Automatic
		common.PrintResult(security.CheckIRSAAndPodIdentity(k8sClient))

		// 데이터 플레인 노드에 필수로 필요한 IAM 권한만 부여 - Automatic
		common.PrintResult(security.CheckNodeIAMRoles(k8sClient))

		// 루트 유저가 아닌 유저로 컨테이너 실행 - Automatic
		common.PrintResult(security.CheckContainerExecutionUser(k8sClient))

		// 불필요한 OS 권한 비부여 - Automatic

		// 멀티 태넌시 적용 유무 - Manual
		common.PrintResult(security.CheckMultitenancy(k8sClient, cfg, cluster))

		// Audit 로그 활성화 - Automatic
		common.PrintResult(security.CheckAuditLoggingEnabled(&security.EksCluster{Cluster: eksCluster.Cluster}))

		// 비정상 접근에 대한 알림 설정 - Manual
		common.PrintResult(security.CheckAccessAlarm())

		// Pod-to-Pod 접근 제어 - Automatic/Manual
		common.PrintResult(security.CheckPodToPodNetworkPolicy(k8sClient, cluster))

		// PV 암호화 - Automatic
		common.PrintResult(security.CheckPVEcryption(k8sClient))

		// Secret 객체 암호화 - Automatic
		common.PrintResult(security.CheckSecretEncryption(k8sClient))

		// 데이터 플레인 사설망 - Automatic
		common.PrintResult(security.DataplanePrivateCheck(security.EksCluster(eksCluster), cfg))

		// 컨테이너 이미지 정적 분석 - Manual
		common.PrintResult(security.CheckImageStaticAnalysis(k8sClient, cfg, cluster))

		// 읽기 전용 파일시스템 사용 - Automatic
		common.PrintResult(security.ReadnonlyFilesystemCheck(k8sClient))

		// Scalability 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Scalability Check]===============\n")

		// Karpenter 사용 - Automatic
		common.PrintResult(scalability.GetKarpenter(k8sClient))

		// Karpenter 전용 노드 그룹 혹은 Fargate 사용 - Automatic
		common.PrintResult(scalability.CheckNodeGroupUsage(k8sClient))

		// Spot 노드 사용시 Spot 중지 핸들러 적용 - Automatic
		common.PrintResult(scalability.CheckSpotNodeTerminationHandler(k8sClient))

		// 중요 Pod에 노드 삭제 방지용 Label 부여 - Manual
		common.PrintResult(scalability.CheckImportantPodProtection(k8sClient, cfg, cluster))

		// Application에 Graceful shutdown 적용 - Manual
		common.PrintResult(scalability.CheckGracefulShutdown(k8sClient, cfg, cluster))

		// 노드 확장/축소 정책 적용 - Manual
		common.PrintResult(scalability.CheckNodeScalingPolicy())

		// 다양한 인스턴스 타입 사용 - Automatic
		common.PrintResult(scalability.CheckInstanceTypes(k8sClient))

		// Scalability 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Stability Check]===============\n")

		// 싱글톤 Pod 미사용 - Automatic
		common.PrintResult(stability.SingletonPodCheck(k8sClient))

		// 2개 이상의 Pod 복제본 사용 - Automatic
		common.PrintResult(stability.PodReplicaSetCheck(k8sClient))

		// 동일한 역할을 하는 Pod를 다수의 노드에 분산 배포 - Automatic
		common.PrintResult(stability.CheckPodDistributionAndAffinity(k8sClient))

		// HPA 적용 - Automatic
		common.PrintResult(stability.CheckHpa(k8sClient))

		// Probe(Startup, Readiness, Liveness) 적용 - Automatic
		common.PrintResult(stability.CheckProbe(k8sClient))

		// 중요 워크로드에 대한 PDB(Pod Distruption Budget) 적용 - Automatic/Manual
		common.PrintResult(stability.CheckPDB())

		// 애플리케이션에 적절한 CPU/RAM 할당 - Automatic/Manual
		common.PrintResult(stability.CheckResourceAllocation(k8sClient, cfg, cluster))

		// 애플리케이션 중요도에 따른 QoS 적용 - Automatic/Manual
		common.PrintResult(stability.CheckQoSClass(k8sClient, cfg, cluster))

		// 오토스케일링 그룹 기반 관리형 노드 그룹 생성 - Automatic
		common.PrintResult(stability.CheckAutoScaledManagedNodeGroup(k8sClient, cluster))

		// Cluster Autoscaler 적용 - Automatic
		common.PrintResult(stability.CheckClusterAutoscalerEnabled(k8sClient))

		// Karpenter 기반 노드 생성 - Automatic
		common.PrintResult(stability.CheckKarpenterNode(scalability.GetKarpenter(k8sClient), dynamicClient))

		// 다수의 가용 영역에 데이터 플레인 노드 배포 - Automatic
		common.PrintResult(stability.CheckNodeMultiAZ(k8sClient))

		// PV 사용시 volume affinity 위반 사항 체크 - Manual (PV 어피니티 전부다 출력)
		common.PrintResult(stability.CheckVolumeAffinity(k8sClient, cfg, cluster))

		// CoreDNS에 HPA 적용 - Automatic
		common.PrintResult(stability.CheckCoreDNSHpa(k8sClient))

		// DNS 캐시 적용 - Automatic
		common.PrintResult(stability.CheckCoreDNSCache(k8sClient))

		// // Karpenter 사용시 DaemonSet에 Priority Class 부여 - Automatic
		common.PrintResult(stability.CheckDaemonSetPriorityClass(scalability.GetKarpenter(k8sClient), k8sClient))

		// Network 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Network Check]===============\n")

		// VPC 서브넷에 충분한 IP 대역대 확보 - Automatic/Manual
		common.PrintResult(network.CheckVpcSubnetIpCapacity(network.EksCluster(eksCluster), cfg))

		// VPC CNI의 Prefix 모드 사용 - Automatic
		common.PrintResult(network.CheckVpcCniPrefixMode(k8sClient))

		// 사용 사례에 맞는 로드밸런서 사용(ALB or NLB) - Manual
		common.PrintResult(network.CheckLoadBalancerUsage(k8sClient, cfg, cluster))

		// AWS Load Balancer Controller 사용 - Automatic
		common.PrintResult(network.CheckAwsLoadBalancerController(k8sClient))

		// ALB/NLB의 대상으로 Pod의 IP 사용 - Automatic
		common.PrintResult(network.CheckAwsLoadBalancerPodIp(k8sClient))

		// Pod Readiness Gate 적용 - Automatic
		common.PrintResult(network.CheckReadinessGateEnabled(k8sClient))

		// kube-proxy에 IPVS 모드 적용 - Automatic
		common.PrintResult(network.CheckKubeProxyIPVSMode(k8sClient))

		// Endpoint 대신 EndpointSlices 사용 - Automatic
		common.PrintResult(network.EndpointSlicesCheck(k8sClient))

		// 비용최적화 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Cost-Optimized Check]===============\n")

		// EKS용 Kubecost 설치 - Automatic
		common.PrintResult(cost.GetKubecost(k8sClient))

		// 요약본
		common.PrintSummary()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&kubeconfigPath, "kubeconfig", "", "Path to the kubeconfig file to use for CLI requests")
	rootCmd.PersistentFlags().StringVar(&kubeconfigContext, "context", "", "The name of the kubeconfig context to use")
}

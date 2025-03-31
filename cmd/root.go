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

const (
	Red    = "\033[31m" // 빨간색
	Green  = "\033[32m" // 초록색
	Yellow = "\033[33m" // 노란색
	Reset  = "\033[0m"  // 기본 색상으로 리셋
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

		// 컨테이너 이미지 태그에 latest 미사용
		common.PrintResult(general.CheckImageTag(k8sClient))

		// Security 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Security Check]===============\n")

		// EKS 클러스터 API 엔드포인트 접근 제어(공인망, 사설망, IP 기반 제어) - Automatic
		common.PrintResult(security.CheckEndpointPublicAccess(security.EksCluster(eksCluster)))

		// 클러스터 접근 제어(Access entries, aws-auth 컨피그맵) - Automatic/Manual
		// 컨피그맵이랑 accesslist 출력인데 정확히 어케 출력되야되는지랑, 인자로 cluster 받는거 맞는지 확인 필요
		common.PrintResult(security.CheckAccessControl(k8sClient, cluster))

		// IRSA 또는 Pod Identity 기반 권한 부여 - Automatic
		common.PrintResult(security.CheckIRSAAndPodIdentity(k8sClient))

		// 데이터 플레인 노드에 필수로 필요한 IAM 권한만 부여 - Automatic
		common.PrintResult(security.CheckNodeIAMRoles(k8sClient))

		// 루트 유저가 아닌 유저로 컨테이너 실행 - Automatic
		common.PrintResult(security.CheckContainerExecutionUser(k8sClient))

		// Audit 로그 활성화 - Automatic
		common.PrintResult(security.CheckAuditLoggingEnabled(&security.EksCluster{Cluster: eksCluster.Cluster}))

		// Manual 포함된 기능이라 출력 템플릿 스킵 - 검토 필요
		// Pod-to-Pod 접근 제어 - Automatic/Manual
		if security.CheckPodToPodNetworkPolicy(k8sClient) {
			fmt.Println(Green + "✔ PASS: Pod-to-Pod network policy is found" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Pod-to-Pod network policy is not found" + Reset)
		}

		// PV 암호화 - Automatic
		common.PrintResult(security.CheckPVEcryption(k8sClient))

		// Secret 객체 암호화 - Automatic
		common.PrintResult(security.CheckSecretEncryption(k8sClient))

		// 데이터 플레인 사설망 - Automatic
		common.PrintResult(security.DataplanePrivateCheck(security.EksCluster(eksCluster), cfg))

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

		// 오토스케일링 그룹 기반 관리형 노드 그룹 생성 - Automatic
		common.PrintResult(stability.CheckAutoScaledManagedNodeGroup(k8sClient, cluster))

		// Cluster Autoscaler 적용 - Automatic
		common.PrintResult(stability.CheckClusterAutoscalerEnabled(k8sClient))

		// Karpenter 기반 노드 생성 - Automatic
		common.PrintResult(stability.CheckKarpenterNode(scalability.GetKarpenter(k8sClient), dynamicClient))

		// 다수의 가용 영역에 데이터 플레인 노드 배포 - Automatic
		common.PrintResult(stability.CheckNodeMultiAZ(k8sClient))

		// CoreDNS에 HPA 적용 - Automatic
		common.PrintResult(stability.CheckCoreDNSHpa(k8sClient))

		// DNS 캐시 적용 - Automatic
		common.PrintResult(stability.CheckCoreDNSCache(k8sClient))

		// // Karpenter 사용시 DaemonSet에 Priority Class 부여 - Automatic
		common.PrintResult(stability.CheckDaemonSetPriorityClass(scalability.GetKarpenter(k8sClient), k8sClient))

		// Network 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Network Check]===============\n" + Reset)

		// VPC 서브넷에 충분한 IP 대역대 확보 - Automatic/Manual
		if ipCapacities := network.CheckVpcSubnetIpCapacity(network.EksCluster(eksCluster), cfg); len(ipCapacities) > 0 {
			for subnetId, ipCapacity := range ipCapacities {
				fmt.Printf(Red+"✖ FAIL: Subnet %s has less than 10%% of available IPs remaining: %d\n", subnetId, ipCapacity)
			}
		} else {
			fmt.Println(Green + "✔ PASS: All subnets have more than 10% of available IPs remaining" + Reset)
		}

		// VPC CNI의 Prefix 모드 사용 - Automatic
		if network.CheckVpcCniPrefixMode(k8sClient) {
			fmt.Println(Green + "✔ PASS: VPC CNI is in prefix mode" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: VPC CNI is not in prefix mode" + Reset)
		}

		// AWS Load Balancer Controller 사용 - Automatic
		if network.CheckAwsLoadBalancerController(k8sClient) {
			fmt.Println(Green + "✔ PASS: AWS Load Balancer Controller is installed" + Reset)
			// ALB/NLB의 대상으로 Pod의 IP 사용 - Automatic
			if network.CheckAwsLoadBalancerPodIp(k8sClient) {
				fmt.Println(Green + "✔ PASS: AWS Load Balancer Pod IP is in use" + Reset)
			} else {
				fmt.Println(Red + "FAIL: AWS Load Balancer Pod IP is not in use" + Reset)
			}
			// Pod Readiness Gate 적용 - Automatic
			if network.CheckReadinessGateEnabled(k8sClient) {
				fmt.Println(Green + "✔ PASS: Readiness Gate is enabled" + Reset)
			} else {
				fmt.Println(Red + "✖ FAIL: Readiness Gate is not enabled" + Reset)
			}
		} else {
			fmt.Println(Red + "✖ FAIL: AWS Load Balancer Controller is not installed" + Reset)
			fmt.Println(Red + "✖ FAIL: AWS Load Balancer Pod IP is not in use" + Reset)
			fmt.Println(Red + "✖ FAIL: Readiness Gate is not enabled" + Reset)
		}

		// kube-proxy에 IPVS 모드 적용 - Automatic
		if network.CheckKubeProxyIPVSMode(k8sClient) {
			fmt.Println(Green + "✔ PASS: kube-proxy is in IPVS mode" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: kube-proxy is not in IPVS mode" + Reset)
		}

		// Endpoint 대신 EndpointSlices 사용 - Automatic
		network.EndpointSlicesCheck(k8sClient)

		// 비용최적화 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Cost-Optimized Check]===============\n")

		// EKS용 Kubecost 설치 - Automatic
		if cost.GetKubecost(k8sClient) {
			fmt.Println(Green + "✔ PASS: Kubecost is installed" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Kubecost is not installed" + Reset)
		}
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

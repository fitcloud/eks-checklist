package cmd

import (
	"fmt"
	"os"

	"eks-checklist/cmd/cost"
	"eks-checklist/cmd/general"
	"eks-checklist/cmd/network"
	"eks-checklist/cmd/scalability"
	"eks-checklist/cmd/security"
	"eks-checklist/cmd/stability"

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

		// latest 태그를 가진 이미지를 사용해서는 안됨
		general.CheckImageTag(k8sClient)

		// Security 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Security Check]===============\n")

		// EKS 클러스터 API 엔드포인트 접근 제어(공인망, 사설망, IP 기반 제어) - Automatic
		if !eksCluster.Cluster.ResourcesVpcConfig.EndpointPublicAccess {
			fmt.Println(Green + "✔ PASS: EKS Cluster is not publicly accessible from the internet" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: EKS Cluster is publicly accessible from the internet" + Reset)
		}

		// 클러스터 접근 제어(Access entries, aws-auth 컨피그맵) - Automatic/Manual
		// 컨피그맵이랑 accesslist 출력인데 정확히 어케 출력되야되는지랑, 인자로 cluster 받는거 맞는지 확인 필요
		security.PrintAccessControl(k8sClient, cluster)

		// 데이터 플레인 노드에 필수로 필요한 IAM 권한만 부여 - Automatic
		security.CheckNodeIAMRoles(k8sClient)

		// 루트 유저가 아닌 유저로 컨테이너 실행 - Automatic
		security.CheckContainerExecutionUser(k8sClient)

		// Audit 로그 활성화 - Automatic
		if security.CheckAuditLoggingEnabled(&security.EksCluster{Cluster: eksCluster.Cluster}) {
			fmt.Println(Green + "✔ PASS: Audit logging is enabled" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Audit logging is not enabled" + Reset)
		}

		// Pod-to-Pod 접근 제어 - Automatic/Manual
		if security.CheckPodToPodNetworkPolicy(k8sClient) {
			fmt.Println(Green + "✔ PASS: Pod-to-Pod network policy is found" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Pod-to-Pod network policy is not found" + Reset)
		}

		// PV 암호화 - Automatic
		if security.CheckPVEcryption(k8sClient) {
			fmt.Println(Green + "✔ PASS: PV encryption is enabled" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: PV encryption is not enabled" + Reset)
		}

		// Secret 객체 암호화 - Automatic
		if security.CheckSecretEncryption(k8sClient) {
			fmt.Println(Green + "✔ PASS: Secret encryption is enabled" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Secret encryption is not enabled" + Reset)
		}

		// 읽기 전용 파일시스템 사용 - Automatic
		security.ReadnonlyFilesystemCheck(k8sClient)

		security.CheckIRSAAndPodIdentity(k8sClient)

		// Scalability 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Scalability Check]===============\n")

		// Karpenter 사용 - Automatic
		if getKarpenter(k8sClient) {
			fmt.Println(Green + "✔ PASS: Karpenter is installed" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Karpenter is not installed" + Reset)
		}

		// Karpenter 전용 노드 그룹 혹은 Fargate 사용 - Automatic
		scalability.CheckNodeGroupUsage(k8sClient)

		// Spot 노드 사용시 Spot 중지 핸들러 적용 - Automatic
		if scalability.CheckSpotNodeTerminationHandler(k8sClient) {
			fmt.Println(Green + "✔ PASS: Spot Node Termination Handler is applied" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Spot Node Termination Handler is not applied" + Reset)
		}

		// 다양한 인스턴스 타입 사용 - Automatic
		if scalability.CheckInstanceTypes(k8sClient) {
			fmt.Println(Green + "✔ PASS: Cluster has multiple instance types" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Cluster does not have multiple instance types" + Reset)
		}

		// Scalability 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Stability Check]===============\n")

		// 싱글톤 Pod 미사용 - Automatic
		if stability.SingletonPodCheck(k8sClient) {
			fmt.Println(Green + "✔ PASS: SingletonPod No Used" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: SingletonPod Used" + Reset)
		}

		// 2개 이상의 Pod 복제본 사용 - Automatic
		if stability.PodReplicaSetCheck(k8sClient) {
			fmt.Println(Green + "✔ PASS: ReplicaSet Used more than one Pod" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: ReplicaSet Used one Pod" + Reset)
		}

		// HPA 적용 - Automatic
		stability.CheckHpa(k8sClient)

		// Probe(Startup, Readiness, Liveness) 적용 - Automatic
		if stability.CheckProbe(k8sClient) {
			fmt.Println(Green + "✔ PASS: Probe is applied" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Probe is not applied" + Reset)
		}

		// 오토스케일링 그룹 기반 관리형 노드 그룹 생성 - Automatic
		if stability.CheckAutoScaledManagedNodeGroup(k8sClient, cluster) {
			fmt.Println(Green + "✔ PASS: AutoScaled Managed Node Group is created" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: AutoScaled Managed Node Group is not created" + Reset)
		}

		// Cluster Autoscaler 적용 - Automatic
		if stability.CheckClusterAutoscalerEnabled(k8sClient) {
			fmt.Println(Green + "✔ PASS: Cluster Autoscaler is installed" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Cluster Autoscaler is not installed" + Reset)
		}

		// Karpenter 기반 노드 생성 - Automatic
		if getKarpenter(k8sClient) {
			if stability.CheckKarpenterNode(dynamicClient) {
				fmt.Println(Green + "✔ PASS: Karpenter Node is created" + Reset)
			} else {
				fmt.Println(Red + "✖ FAIL: Karpenter Node is not created" + Reset)
			}
		} else {
			fmt.Println(Yellow + "⚠ WARNING: Karpenter is not installed" + Reset)
		}

		// 다수의 가용 영역에 데이터 플레인 노드 배포 - Automatic
		if stability.CheckNodeMultiAZ(k8sClient) {
			fmt.Println(Green + "✔ PASS: Nodes are deployed across multiple availability zones" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Nodes are not deployed across multiple availability zones" + Reset)
		}

		// CoreDNS의 HPA가 존재하는지 확인
		if stability.CheckCoreDNSHpa(k8sClient) {
			fmt.Println(Green + "✔ PASS: CoreDNS HPA is installed" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: CoreDNS HPA is not installed" + Reset)
		}

		// DNS 캐시 적용 - Automatic
		if stability.CheckCoreDNSCache(k8sClient) {
			fmt.Println(Green + "✔ PASS: CoreDNS Cache is applied" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: CoreDNS Cache is not applied" + Reset)
		}

		// Karpenter 사용시 DaemonSet에 Priority Class 부여 - Automatic
		if getKarpenter(k8sClient) {
			if stability.CheckDaemonSetPriorityClass(k8sClient) {
				fmt.Println(Green + "✔ PASS: DaemonSet Priority Class is set" + Reset)
			} else {
				fmt.Println(Red + "✖ FAIL: DaemonSet Priority Class is not set" + Reset)
			}
		} else {
			fmt.Println(Yellow + "⚠ WARNING: Karpenter is not installed" + Reset)
		}

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

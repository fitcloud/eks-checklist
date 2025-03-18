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

		// General 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[General Check]===============\n")

		// latest 태그를 가진 이미지를 사용해서는 안됨
		general.CheckImageTag(k8sClient)

		// Security 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Security Check]===============\n")

		// 클러스터 엔드포인트가 public 인지 않인지 확인
		if !eksCluster.Cluster.ResourcesVpcConfig.EndpointPublicAccess {
			fmt.Println(Green + "✔ PASS: EKS Cluster is not publicly accessible from the internet" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: EKS Cluster is publicly accessible from the internet" + Reset)
		}

		// EKS Audit log 활성화 유무 확인
		if security.CheckAuditLoggingEnabled(&security.EksCluster{Cluster: eksCluster.Cluster}) {
			fmt.Println(Green + "✔ PASS: Audit logging is enabled" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Audit logging is not enabled" + Reset)
		}

		// 컨테이너 이미지들이 root을 사용하면 안됨
		security.CheckContainerExecutionUser(k8sClient)

		// worker node role에는 최소환의 권한만 가져야함
		security.CheckNodeIAMRoles(k8sClient)

		// Scalability 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Scalability Check]===============\n")

		// Karpenter 사용 여부 확인
		if getKarpenter(k8sClient) {
			fmt.Println(Green + "✔ PASS: Karpenter is installed" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Karpenter is not installed" + Reset)
		}

		// 클러스터의 노드 인스턴스 유형이 다양한지 확인
		if scalability.CheckInstanceTypes(k8sClient) {
			fmt.Println(Green + "✔ PASS: Cluster has multiple instance types" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Cluster does not have multiple instance types" + Reset)
		}

		// Scalability 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Stability Check]===============\n")

		// 클러스터에 Cluster Autoscaler가 설치되어 있는지 확인
		if stability.CheckClusterAutoscalerEnabled(k8sClient) {
			fmt.Println(Green + "✔ PASS: Cluster Autoscaler is installed" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: Cluster Autoscaler is not installed" + Reset)
		}

		// CoreDNS의 HPA가 존재하는지 확인
		if stability.CheckCoreDNSHpa(k8sClient) {
			fmt.Println(Green + "✔ PASS: CoreDNS HPA is installed" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: CoreDNS HPA is not installed" + Reset)
		}

		// 싱글톤 Pod 사용 중인지 확인
		if stability.SingletonPodCheck(k8sClient) {
			fmt.Println(Green + "✔ PASS: SingletonPod No Used" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: SingletonPod Used" + Reset)
		}

		if stability.PodReplicaSetCheck(k8sClient) {
			fmt.Println(Green + "✔ PASS: ReplicaSet Used more than one Pod" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: ReplicaSet Used one Pod" + Reset)
		}

		// Network 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Network Check]===============\n" + Reset)

		// 클러스터가 사용하는 서브넷의 가용가능한 IP주소의 가용성 검사하는 함수, 반환값이 있을 경우 FAIL 출력 없으면 PASS 출력
		if ipCapacities := network.CheckVpcSubnetIpCapacity(network.EksCluster(eksCluster)); len(ipCapacities) > 0 {
			for subnetId, ipCapacity := range ipCapacities {
				fmt.Printf(Red+"✖ FAIL: Subnet %s has less than 10%% of available IPs remaining: %d\n", subnetId, ipCapacity)
			}
		} else {
			// 해석하면 모든 서브넷에 사용가능한 IP주소가 10%이상 남아있다
			fmt.Println(Green + "✔ PASS: All subnets have more than 10% of available IPs remaining" + Reset)
		}

		// VPC CNI에서 Prefix 모드 사용 유무 확인
		if network.CheckVpcCniPrefixMode(k8sClient) {
			fmt.Println(Green + "✔ PASS: VPC CNI is in prefix mode" + Reset)
		} else {
			fmt.Println(Red + "✖ FAIL: VPC CNI is not in prefix mode" + Reset)
		}

		// aws-loadblaancer-controller 설치 여부 확인
		if network.CheckAwsLoadBalancerController(k8sClient) {
			fmt.Println(Green + "✔ PASS: AWS Load Balancer Controller is installed" + Reset)
			// AWS Load Balancer Pod IP가 사용 중인지 확인
			if network.CheckAwsLoadBalancerPodIp(k8sClient) {
				fmt.Println(Green + "✔ PASS: AWS Load Balancer Pod IP is in use" + Reset)
			} else {
				fmt.Println(Red + "FAIL: AWS Load Balancer Pod IP is not in use" + Reset)
			}
			// Readiness Gate 활성화 유무
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

		// 클러스터에 Horizontal Pod Autoscaler가 설정되어 있는지 확인
		stability.CheckHpa(k8sClient)

		// 비용최적화 항목 체크 기능은 하단 항목에 추가
		fmt.Printf("\n===============[Cost-Optimized Check]===============\n")

		// 클러스터에 Kubecost가 설치되어 있는지 확인
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

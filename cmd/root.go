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

var (
	kubeconfigPath    string
	kubeconfigContext string
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "eks-checklist",
	Long:  "eks-checklist",
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig := getKubeconfig(kubeconfigPath)
		cluster := getEksClusterName(kubeconfig)

		fmt.Printf("Running checks on %s\n", cluster)

		eksCluster := Describe(cluster)
		k8sClient := createK8sClient(kubeconfig)

		// 클러스터 엔드포인트가 public 인지 않인지 확인
		if !eksCluster.Cluster.ResourcesVpcConfig.EndpointPublicAccess {
			fmt.Println("PASS: EKS Cluster is not publicly accessible from the internet")
		} else {
			fmt.Println("FAIL: EKS Cluster is publicly accessible from the internet")
		}

		// 클러스터가 사용하는 서브넷의 가용가능한 IP주소의 가용성 검사하는 함수, 반환값이 있을 경우 FAIL 출력 없으면 PASS 출력
		if ipCapacities := network.CheckVpcSubnetIpCapacity(network.EksCluster(eksCluster)); len(ipCapacities) > 0 {
			for subnetId, ipCapacity := range ipCapacities {
				fmt.Printf("FAIL: Subnet %s has less than 10%% of available IPs remaining: %d\n", subnetId, ipCapacity)
			}
		} else {
			// 해석하면 모든 서브넷에 사용가능한 IP주소가 10%이상 남아있다
			fmt.Println("PASS: All subnets have more than 10% of available IPs remaining")
		}

		if getKarpenter(k8sClient) {
			fmt.Println("PASS: Karpenter is installed")
		} else {
			fmt.Println("FAIL: Karpenter is not installed")
		}

		// VPC CNI에서 Prefix 모드 사용 유무 확인
		if network.CheckVpcCniPrefixMode(k8sClient) {
			fmt.Println("PASS: VPC CNI is in prefix mode")
		} else {
			fmt.Println("FAIL: VPC CNI is not in prefix mode")
		}

		// EKS Audit log 활성화 유무 확인
		if security.CheckAuditLoggingEnabled(&security.EksCluster{Cluster: eksCluster.Cluster}) {
			fmt.Println("PASS: Audit logging is enabled")
		} else {
			fmt.Println("FAIL: Audit logging is not enabled")
		}

		// aws-loadblaancer-controller 설치 여부 확인
		if network.CheckAwsLoadBalancerController(k8sClient) {
			fmt.Println("PASS: AWS Load Balancer Controller is installed")
			// AWS Load Balancer Pod IP가 사용 중인지 확인
			if network.CheckAwsLoadBalancerPodIp(k8sClient) {
				fmt.Println("PASS: AWS Load Balancer Pod IP is in use")
			} else {
				fmt.Println("FAIL: AWS Load Balancer Pod IP is not in use")
			}
		} else {
			fmt.Println("FAIL: AWS Load Balancer Controller is not installed")
			fmt.Println("FAIL: AWS Load Balancer Pod IP is not in use")
		}

		// latest 태그를 가진 이미지를 사용해서는 안됨
		general.CheckImageTag(k8sClient)

		// 컨테이너 이미지들이 root을 사용하면 안됨
		security.CheckContainerExecutionUser(k8sClient)

		// 클러스터에 Cluster Autoscaler가 설치되어 있는지 확인
		if stability.CheckClusterAutoscalerEnabled(k8sClient) {
			fmt.Println("PASS: Cluster Autoscaler is installed")
		} else {
			fmt.Println("FAIL: Cluster Autoscaler is not installed")
		}

		// 클러스터의 노드 인스턴스 유형이 다양한지 확인
		if scalability.CheckInstanceTypes(k8sClient) {
			fmt.Println("PASS: Cluster has multiple instance types")
		} else {
			fmt.Println("FAIL: Cluster does not have multiple instance types")
		}

		// 클러스터에 Kubecost가 설치되어 있는지 확인
		if cost.GetKubecost(k8sClient) {
			fmt.Println("PASS: Kubecost is installed")
		} else {
			fmt.Println("FAIL: Kubecost is not installed")
		}

		// 싱글톤 Pod 사용 중인지 확인
		if stability.SingletonPodCheck(k8sClient) {
			fmt.Println("PASS: SingletonPod No Used")
		} else {
			fmt.Println("FAIL: SingletonPod Used")
		}

		// CoreDNS의 HPA가 존재하는지 확인
		if stability.CheckCoreDNSHpa(k8sClient) {
			fmt.Println("PASS: CoreDNS HPA is installed")
		} else {
			fmt.Println("FAIL: CoreDNS HPA is not installed")
		}

		// worker node role에는 최소환의 권한만 가져야함
		security.CheckNodeIAMRoles(k8sClient)

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

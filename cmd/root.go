package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	// 승도가 만든 네트워크 패키지
	"eks-checklist/cmd/general"
	"eks-checklist/cmd/network"
	"eks-checklist/cmd/security"
	"eks-checklist/cmd/stability"
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

		k8sClient := createK8sClient(kubeconfig)

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

		// EKS 클러스터 Audit로깅 설정 유무를 확인하는 함수, 존재하면 PASS 출력 아니면 FAIL 출력
		if eksCluster.Cluster.Logging.ClusterLogging[0].Types[0] == "audit" {
			fmt.Println("PASS: Audit logging is enabled")
		} else {
			fmt.Println("FAIL: Audit logging is not enabled")
		}

		// aws-loadblaancer-controller 설치 여부 확인
		if network.CheckAwsLoadBalancerController(k8sClient) {
			fmt.Println("PASS: AWS Load Balancer Controller is installed")
		} else {
			fmt.Println("FAIL: AWS Load Balancer Controller is not installed")
		}

		// AWS Load Balancer Pod IP가 사용 중인지 확인
		if network.CheckAwsLoadBalancerPodIp(k8sClient) {
			fmt.Println("PASS: AWS Load Balancer Pod IP is in use")
		} else {
			fmt.Println("FAIL: AWS Load Balancer Pod IP is not in use")
		}

		// 이렇게 그냥 함수 안에 if로 넣어도 될까용?
		general.CheckImageTag(k8sClient)
		security.CheckContainerExecutionUser(k8sClient)

		// 클러스터에 Cluster Autoscaler가 설치되어 있는지 확인
		if stability.CheckClusterAutoscalerEnabled(k8sClient) {
			fmt.Println("PASS: Cluster Autoscaler is installed")
		} else {
			fmt.Println("FAIL: Cluster Autoscaler is not installed")
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

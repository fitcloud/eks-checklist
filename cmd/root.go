package cmd

import (
	"fmt"
	"os"

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

		if !eksCluster.Cluster.ResourcesVpcConfig.EndpointPublicAccess {
			fmt.Println(C("PASS: EKS Cluster is not publicly accessible from the internet"))
		} else {
			fmt.Println(C("FAIL: EKS Cluster is publicly accessible from the internet"))
		}

		if eksCluster.Cluster.Logging.ClusterLogging[0].Types[0] == "audit" {
			fmt.Println(C("PASS: Audit logging is enabled"))
		} else {
			fmt.Println(C("FAIL: Audit logging is not enabled"))
		}

		if eksCluster.Cluster.AccessConfig.AuthenticationMode == "API_AND_CONFIG_MAP" {
			fmt.Println(C("PASS: RBAC is enabled"))
		} else {
			fmt.Println(C("FAIL: RBAC is not enabled"))
		}

		if eksCluster.Cluster.UpgradePolicy.SupportType == "EXTENDED" {
			fmt.Println(C("PASS: Extended support is enabled"))
		} else {
			fmt.Println(C("FAIL: Extended support is not enabled"))
		}

		k8sClient := createK8sClient(kubeconfig)

		if getKarpenter(k8sClient) {
			fmt.Println(C("PASS: Karpenter is installed"))
		} else {
			fmt.Println(C("FAIL: Karpenter is not installed"))
		}

		if getClusterAutoscaler(k8sClient) {
			fmt.Println(C("PASS: Cluster Autoscaler is installed"))
		} else {
			fmt.Println(C("FAIL: Cluster Autoscaler is not installed"))
		}

		if getKubecost(k8sClient) {
			fmt.Println(C("PASS: Kubecost is installed"))
		} else {
			fmt.Println(C("FAIL: Kubecost is not installed"))
		}

		if getHpa(k8sClient) {
			fmt.Println(C("PASS: Horizontal Pod Autoscaler is installed"))
		} else {
			fmt.Println(C("FAIL: Horizontal Pod Autoscaler is not installed"))
		}

		if getcontainerImagelatestTag(k8sClient) {
			fmt.Println(C("PASS: Containers are not using the latest tag"))
		} else {
			fmt.Println(C("FAIL: Containers are using the latest tag"))
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

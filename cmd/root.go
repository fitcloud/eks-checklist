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
			fmt.Println("PASS: EKS Cluster is not publicly accessible from the internet")
		} else {
			fmt.Println("FAIL: EKS Cluster is publicly accessible from the internet")
		}

		k8sClient := createK8sClient(kubeconfig)

		if getKarpenter(k8sClient) {
			fmt.Println("PASS: Karpenter is installed")
		} else {
			fmt.Println("FAIL: Karpenter is not installed")
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

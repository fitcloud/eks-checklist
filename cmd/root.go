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
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

var (
	kubeconfigPath    string
	kubeconfigContext string
	awsProfile        string
	outputFilter      string
	outputFormat      string
	sortMode          bool
)

var rootCmd = &cobra.Command{
	Use:   os.Args[0],
	Short: "eks-checklist",
	Long:  "eks-checklist",
	Run: func(cmd *cobra.Command, args []string) {
		// 설정 초기화
		setupConfig()
		
		// 클러스터 및 클라이언트 초기화
		AWS_PROFILE, kubeconfig, err := getKubeconfig(kubeconfigPath, kubeconfigContext, awsProfile)
		if err != nil {
			fmt.Printf("오류: %v\n", err)
			os.Exit(1)
		}

		cluster := getEksClusterName(kubeconfig)
		fmt.Printf("Running checks on %s\n", cluster)
		
		k8sClient, err := createK8sClient(kubeconfig)
		if err != nil {
			fmt.Printf("오류: %v\n", err)
			os.Exit(1)
		}

		cfg, err := GetAWSConfig(AWS_PROFILE)
		if err != nil {
			fmt.Printf("오류: %v\n", err)
			os.Exit(1)
		}

		eksCluster := Describe(cluster, cfg)
		dynamicClient, err := CreateDynamicClient(&kubeconfig)
		if err != nil {
			fmt.Printf("오류: %v\n", err)
			os.Exit(1)
		}

		// 체크 레지스트리 생성
		registry := common.NewCheckerRegistry()

		// 카테고리별 체크 항목 등록
		general.RegisterCheckers(registry, k8sClient)
		security.RegisterCheckers(registry, k8sClient, cfg, cluster, eksCluster)
		scalability.RegisterCheckers(registry, k8sClient, cfg, cluster)
		stability.RegisterCheckers(registry, k8sClient, cfg, cluster, dynamicClient)
		network.RegisterCheckers(registry, k8sClient, cfg, cluster)
		cost.RegisterCheckers(registry, k8sClient, cfg, cluster)

		// 모든 체크 실행
		registry.RunChecks()
	},
}

// 설정 초기화 함수
func setupConfig() {
	common.SetSortMode(sortMode)

	if outputFilter != "" {
		// 소문자로 변환하여 비교
		lowerFilter := strings.ToLower(outputFilter)
		validFilters := common.GetValidOutputFilters()
		isValid := false

		for _, valid := range validFilters {
			if lowerFilter == valid {
				isValid = true
				break
			}
		}

		if !isValid {
			fmt.Printf("오류: 유효하지 않은 출력 필터 '%s'\n", outputFilter)
			fmt.Printf("유효한 값: %s\n", strings.Join(validFilters, ", "))
			os.Exit(1)
		}

		fmt.Printf("Output filter: %s\n", lowerFilter)
		common.SetOutputFilter(lowerFilter)
	}

	// 출력 형식 설정
	if outputFormat != "" {
		lowerFormat := strings.ToLower(outputFormat)
		validFormats := common.GetValidOutputFormats()
		isValid := false

		for _, valid := range validFormats {
			if lowerFormat == valid {
				isValid = true
				break
			}
		}

		if !isValid {
			fmt.Printf("오류: 유효하지 않은 출력 형식 '%s'\n", outputFormat)
			fmt.Printf("유효한 값: %s\n", strings.Join(validFormats, ", "))
			os.Exit(1)
		}

		common.SetOutputFormat(lowerFormat)
	}

	// HTML 출력 초기화
	if outputFormat == common.OutputFormatHTML || outputFormat == common.OutputFormatPDF {
		common.InitHTMLOutput()
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	home := homedir.HomeDir()
	defaultKubeconfigPath := filepath.Join(home, ".kube", "config")

	rootCmd.PersistentFlags().StringVar(&kubeconfigPath, "kubeconfig", defaultKubeconfigPath, "kubeconfig 파일 경로")
	rootCmd.PersistentFlags().StringVar(&kubeconfigContext, "context", "", "사용할 kubeconfig 컨텍스트")
	rootCmd.PersistentFlags().StringVar(&awsProfile, "profile", "", "사용할 AWS 프로파일")
	rootCmd.PersistentFlags().StringVar(&outputFilter, "filter", "", "출력 필터 (all, pass, fail, manual)")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "output", "text", "출력 형식 (text, html, pdf)")
	rootCmd.PersistentFlags().BoolVar(&sortMode, "sort", false, "결과를 상태별로 정렬")
}

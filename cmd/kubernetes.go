// package cmd

// import (
// 	// "flag"

// 	"os"

// 	// "path/filepath"
// 	"slices"

// 	"k8s.io/client-go/dynamic"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/rest"
// 	"k8s.io/client-go/tools/clientcmd"
// 	// "k8s.io/client-go/util/homedir"
// )

// func getKubeconfig(kubeconfigPath string, awsProfile string) rest.Config {
// 	kubeconfig := &kubeconfigPath

// 	// if home := homedir.HomeDir(); kubeconfigPath == "" {
// 	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
// 	// } else {
// 	// 	kubeconfig = flag.String("kubeconfig", kubeconfigPath, "absolute path to the kubeconfig file")
// 	// }

// 	// flag.Parse()

// 	// AWS_PROFILE 설정
// 	if awsProfile != "" {
// 		os.Setenv("AWS_PROFILE", awsProfile)
// 	}

// 	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	return *config
// }

// func getEksClusterName(kubeconfig rest.Config) string {
// 	clusterNameIdx := slices.Index(kubeconfig.ExecProvider.Args, "--cluster-name") + 1

// 	return kubeconfig.ExecProvider.Args[clusterNameIdx]
// }

// func createK8sClient(kubeconfig rest.Config) kubernetes.Interface {
// 	client, err := kubernetes.NewForConfig(&kubeconfig)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	return client
// }

// // CreateDynamicClient: dynamic.Interface 생성
// func CreateDynamicClient(kubeconfig *rest.Config) (dynamic.Interface, error) {
// 	dynamicClient, err := dynamic.NewForConfig(kubeconfig)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	return dynamicClient, nil
// }

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// getAvailableClusters kubeconfig에서 사용 가능한 모든 클러스터를 가져옵니다
func getAvailableClusters(kubeconfigPath string) ([]string, map[string]string, error) {
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return nil, nil, fmt.Errorf("kubeconfig 파일을 로드하는 중 오류 발생: %v", err)
	}

	// EKS 클러스터만 필터링하기 위한 맵 (컨텍스트 이름 -> 클러스터 이름)
	contextToCluster := make(map[string]string)
	var eksContexts []string

	// 모든 컨텍스트를 확인하여 EKS 클러스터 필터링
	for contextName, context := range config.Contexts {
		clusterName := context.Cluster
		authInfoName := context.AuthInfo

		// 인증 정보 확인
		if authInfo, ok := config.AuthInfos[authInfoName]; ok {
			// EKS 클러스터 식별 (AWS eks get-token이 포함된 경우)
			if authInfo.Exec != nil &&
				authInfo.Exec.Command == "aws" &&
				slices.Contains(authInfo.Exec.Args, "eks") &&
				slices.Contains(authInfo.Exec.Args, "get-token") {
				eksContexts = append(eksContexts, contextName)
				contextToCluster[contextName] = clusterName
			}
		}
	}

	return eksContexts, contextToCluster, nil
}

// selectCluster 사용자에게 클러스터 선택을 위한 대화형 메뉴를 표시합니다
func selectCluster(kubeconfigPath string) (string, error) {
	clusters, contextToCluster, err := getAvailableClusters(kubeconfigPath)
	if err != nil {
		return "", err
	}

	if len(clusters) == 0 {
		return "", fmt.Errorf("kubeconfig에 사용 가능한 EKS 클러스터가 없습니다")
	}

	fmt.Println("\n사용 가능한 EKS 클러스터:")
	for i, cluster := range clusters {
		fmt.Printf("[%d] %s (클러스터: %s)\n", i+1, cluster, contextToCluster[cluster])
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n클러스터 번호를 선택하세요: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("입력 오류:", err)
			continue
		}

		input = strings.TrimSpace(input)
		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(clusters) {
			fmt.Printf("올바른 번호를 입력하세요 (1-%d)\n", len(clusters))
			continue
		}

		selectedContext := clusters[index-1]
		fmt.Printf("\n선택된 클러스터: %s\n", selectedContext)
		return selectedContext, nil
	}
}

// getKubeconfigWithContext는 kubeconfig 파일과 선택된 컨텍스트를 사용하여 rest.Config를 생성합니다
func getKubeconfigWithContext(kubeconfigPath string, context string, awsProfile string) (*rest.Config, error) {
	// AWS_PROFILE 설정
	if awsProfile != "" {
		os.Setenv("AWS_PROFILE", awsProfile)
	}

	// clientcmd.BuildConfigFromFlags는 context 매개변수가 비어 있지 않으면 지정된 컨텍스트를 사용합니다
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{CurrentContext: context},
	).ClientConfig()

	if err != nil {
		return nil, fmt.Errorf("kubeconfig 설정을 로드하는 중 오류 발생: %v", err)
	}

	return config, nil
}

// getKubeconfig 클러스터 선택 기능을 통합한 kubeconfig 로드 함수
func getKubeconfig(kubeconfigPath string, kubeconfigContext string, awsProfile string) rest.Config {
	var config *rest.Config
	var err error

	// 컨텍스트가 명시적으로 지정된 경우 해당 컨텍스트 사용
	if kubeconfigContext != "" {
		config, err = getKubeconfigWithContext(kubeconfigPath, kubeconfigContext, awsProfile)
		if err != nil {
			fmt.Printf("지정된 컨텍스트 '%s'를 로드하는 중 오류 발생: %v\n", kubeconfigContext, err)
			os.Exit(1)
		}
	} else {
		// 대화형 선택 메뉴 표시
		selectedContext, err := selectCluster(kubeconfigPath)
		if err != nil {
			fmt.Printf("클러스터 선택 중 오류 발생: %v\n", err)
			os.Exit(1)
		}

		config, err = getKubeconfigWithContext(kubeconfigPath, selectedContext, awsProfile)
		if err != nil {
			fmt.Printf("선택한 컨텍스트 '%s'를 로드하는 중 오류 발생: %v\n", selectedContext, err)
			os.Exit(1)
		}
	}

	return *config
}

func getEksClusterName(kubeconfig rest.Config) string {
	// ExecProvider가 없는 경우 처리
	if kubeconfig.ExecProvider == nil || len(kubeconfig.ExecProvider.Args) == 0 {
		return "unknown-cluster"
	}

	clusterNameIdx := slices.Index(kubeconfig.ExecProvider.Args, "--cluster-name")
	if clusterNameIdx == -1 || clusterNameIdx+1 >= len(kubeconfig.ExecProvider.Args) {
		return "unknown-cluster"
	}

	return kubeconfig.ExecProvider.Args[clusterNameIdx+1]
}

func createK8sClient(kubeconfig rest.Config) kubernetes.Interface {
	client, err := kubernetes.NewForConfig(&kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	return client
}

// CreateDynamicClient: dynamic.Interface 생성
func CreateDynamicClient(kubeconfig *rest.Config) (dynamic.Interface, error) {
	dynamicClient, err := dynamic.NewForConfig(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	return dynamicClient, nil
}

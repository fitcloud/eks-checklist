package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// selectCluster 함수 수정
func selectCluster(kubeconfigPath string) (string, error) {
	clusters, contextToCluster, err := getAvailableClusters(kubeconfigPath)
	if err != nil {
		return "", err
	}

	if len(clusters) == 0 {
		return "", fmt.Errorf("kubeconfig에 사용 가능한 EKS 클러스터가 없습니다")
	}

	// 1. Kubernetes 클러스터 내부인 경우 (처리는 getKubeconfig에서 이미 했으므로 여기서는 처리 안함)
	// 2. Docker 환경이거나 stdin이 tty가 아닌 경우 첫 번째 클러스터 자동 선택
	if (os.Getenv("RUNNING_IN_DOCKER") == "true" && !isTerminal()) ||
		len(clusters) == 1 { // 클러스터가 하나만 있는 경우도 자동 선택
		fmt.Printf("\n비대화형 환경이거나 클러스터가 하나만 있습니다. 첫 번째 클러스터를 자동으로 선택합니다.\n")
		fmt.Printf("선택된 클러스터: %s (클러스터: %s)\n",
			clusters[0], contextToCluster[clusters[0]])
		return clusters[0], nil
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
			if err == io.EOF {
				// EOF 오류는 파이프된 입력에서 발생할 수 있으므로 첫 번째 클러스터 선택
				fmt.Println("입력을 받을 수 없습니다. 첫 번째 클러스터를 자동으로 선택합니다.")
				return clusters[0], nil
			}
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

// isTerminal은 현재 stdin이 터미널인지 확인합니다
func isTerminal() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
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
func getKubeconfig(kubeconfigPath string, kubeconfigContext string, awsProfile string) (string, rest.Config, error) {
	var config *rest.Config
	var err error
	var selectedContext string
	var AWS_PROFILE string

	// Kubernetes 클러스터 내부에서 실행 중인지 확인
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		fmt.Println("Kubernetes 클러스터 내부에서 실행 중입니다. ServiceAccount 인증 정보를 사용합니다.")

		// InClusterConfig 사용
		config, err = rest.InClusterConfig()
		if err != nil {
			return "", rest.Config{}, fmt.Errorf("인클러스터 설정을 로드하는 중 오류 발생: %w", err)
		}

		// 클러스터 내부에서는 AWS_PROFILE이 의미가 없으므로 빈 문자열 반환
		return "", *config, nil
	}

	// Docker 환경에서 실행 중이지만 클러스터 외부인 경우
	// 컨텍스트가 명시적으로 지정된 경우 해당 컨텍스트 사용
	if kubeconfigContext != "" {
		selectedContext = kubeconfigContext
		config, err = getKubeconfigWithContext(kubeconfigPath, selectedContext, awsProfile)
		if err != nil {
			return "", rest.Config{}, fmt.Errorf("지정된 컨텍스트 '%s'를 로드하는 중 오류 발생: %w", selectedContext, err)
		}
	} else {
		// 대화형 선택 메뉴 표시
		var selErr error
		selectedContext, selErr = selectCluster(kubeconfigPath)
		if selErr != nil {
			return "", rest.Config{}, fmt.Errorf("클러스터 선택 중 오류 발생: %w", selErr)
		}

		config, err = getKubeconfigWithContext(kubeconfigPath, selectedContext, awsProfile)
		if err != nil {
			return "", rest.Config{}, fmt.Errorf("선택한 컨텍스트 '%s'를 로드하는 중 오류 발생: %w", selectedContext, err)
		}
	}

	AWS_PROFILE = getAwsProfileFromContext(kubeconfigPath, selectedContext)

	return AWS_PROFILE, *config, nil
}

// getAwsProfileFromContext는 주어진 컨텍스트에서 AWS_PROFILE 환경 변수 값을 추출합니다
func getAwsProfileFromContext(kubeconfigPath string, contextName string) string {
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		fmt.Printf("kubeconfig 파일을 로드하는 중 오류 발생: %v\n", err)
		return ""
	}

	// 컨텍스트가 존재하는지 확인
	context, exists := config.Contexts[contextName]
	if !exists {
		fmt.Printf("지정된 컨텍스트 '%s'를 찾을 수 없습니다\n", contextName)
		return ""
	}

	// 컨텍스트에 연결된 유저 정보 확인
	authInfoName := context.AuthInfo
	authInfo, exists := config.AuthInfos[authInfoName]
	if !exists {
		fmt.Printf("컨텍스트 '%s'의 유저 정보를 찾을 수 없습니다\n", contextName)
		return ""
	}

	// exec 설정이 있는지 확인
	if authInfo.Exec == nil {
		return ""
	}

	// AWS_PROFILE 환경 변수 검색
	for _, env := range authInfo.Exec.Env {
		if env.Name == "AWS_PROFILE" {
			return env.Value
		}
	}

	return ""
}

// getEksClusterName은 실행 환경에 따라 적절한 방법으로 EKS 클러스터 이름을 가져옵니다
func getEksClusterName(kubeconfig rest.Config) string {
	// 클러스터 내부에서 실행 중인지 확인
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		// 클러스터 내부에서는 다양한 방법으로 클러스터 이름 찾기 시도

		// 방법 1: ConfigMap에서 클러스터 정보 읽기
		clientset, err := kubernetes.NewForConfig(&kubeconfig)
		if err == nil {
			// kube-system 네임스페이스의 aws-auth ConfigMap 확인
			configMap, err := clientset.CoreV1().ConfigMaps("kube-system").Get(context.Background(), "aws-auth", metav1.GetOptions{})
			if err == nil && configMap != nil {
				if clusterName, ok := configMap.Data["cluster-name"]; ok {
					return clusterName
				}
			}

			// 또는 EKS 클러스터의 경우 아래와 같은 ConfigMap에서도 정보를 찾을 수 있음
			configMap, err = clientset.CoreV1().ConfigMaps("kube-system").Get(context.Background(), "cluster-info", metav1.GetOptions{})
			if err == nil && configMap != nil {
				if data, ok := configMap.Data["cluster.name"]; ok {
					return data
				}
			}
		}

		// 방법 2: 환경 변수에서 가져오기
		if clusterName := os.Getenv("CLUSTER_NAME"); clusterName != "" {
			return clusterName
		}

		// 방법 3: 노드 이름에서 추출 시도
		// EKS 노드 이름은 보통 ip-xxx-xxx-xxx-xxx.region.compute.internal 형식
		nodeName := os.Getenv("NODE_NAME")
		if strings.Contains(nodeName, "compute.internal") {
			parts := strings.Split(nodeName, ".")
			if len(parts) > 1 {
				return "eks-cluster-in-" + parts[1]
			}
		}

		return "in-cluster"
	}

	// 클러스터 외부에서는 ExecProvider에서 클러스터 이름 추출
	if kubeconfig.ExecProvider == nil || len(kubeconfig.ExecProvider.Args) == 0 {
		return "unknown-cluster"
	}

	clusterNameIdx := slices.Index(kubeconfig.ExecProvider.Args, "--cluster-name")
	if clusterNameIdx == -1 || clusterNameIdx+1 >= len(kubeconfig.ExecProvider.Args) {
		return "unknown-cluster"
	}

	return kubeconfig.ExecProvider.Args[clusterNameIdx+1]
}

func createK8sClient(kubeconfig rest.Config) (kubernetes.Interface, error) {
	client, err := kubernetes.NewForConfig(&kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Kubernetes 클라이언트 생성 실패: %w", err)
	}

	return client, nil
}

// CreateDynamicClient: dynamic.Interface 생성
func CreateDynamicClient(kubeconfig *rest.Config) (dynamic.Interface, error) {
	dynamicClient, err := dynamic.NewForConfig(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Dynamic 클라이언트 생성 실패: %w", err)
	}

	return dynamicClient, nil
}

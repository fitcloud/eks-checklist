package cmd

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
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
func getKubeconfig(kubeconfigPath string, kubeconfigContext string, awsProfile string) (string, rest.Config) {
	var config *rest.Config
	var err error
	var selectedContext string
	var AWS_PROFILE string

	// Kubernetes 클러스터 내부에서 실행 중인지 확인
	if os.Getenv("IN_K8S") != "" {
		fmt.Println("Kubernetes 클러스터 내부에서 실행 중입니다. ServiceAccount 인증 정보를 사용합니다.")

		// InClusterConfig 사용
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("인클러스터 설정을 로드하는 중 오류 발생: %v\n", err)
			os.Exit(1)
		}

		// 클러스터 내부에서는 AWS_PROFILE이 의미가 없으므로 빈 문자열 반환
		return "", *config
	}

	// Docker 환경에서 실행 중이지만 클러스터 외부인 경우
	// 컨텍스트가 명시적으로 지정된 경우 해당 컨텍스트 사용
	if kubeconfigContext != "" {
		selectedContext = kubeconfigContext
		config, err = getKubeconfigWithContext(kubeconfigPath, selectedContext, awsProfile)
		if err != nil {
			fmt.Printf("지정된 컨텍스트 '%s'를 로드하는 중 오류 발생: %v\n", selectedContext, err)
			os.Exit(1)
		}
	} else {
		// 대화형 선택 메뉴 표시
		var selErr error
		selectedContext, selErr = selectCluster(kubeconfigPath)
		if selErr != nil {
			fmt.Printf("클러스터 선택 중 오류 발생: %v\n", selErr)
			os.Exit(1)
		}

		config, err = getKubeconfigWithContext(kubeconfigPath, selectedContext, awsProfile)
		if err != nil {
			fmt.Printf("선택한 컨텍스트 '%s'를 로드하는 중 오류 발생: %v\n", selectedContext, err)
			os.Exit(1)
		}
	}

	AWS_PROFILE = getAwsProfileFromContext(kubeconfigPath, selectedContext)

	return AWS_PROFILE, *config
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

// JWT 토큰의 페이로드 구조체
type TokenPayload struct {
	Iss string `json:"iss"`
	// 다른 필드들...
}

// getEksClusterName 함수 수정 - AWS 설정을 인자로 받도록 변경
func getEksClusterName(kubeconfig rest.Config, cfg aws.Config) string {
	// 클러스터 내부에서 실행 중인지 확인
	if os.Getenv("IN_K8S") != "" {
		fmt.Println("Kubernetes 클러스터 내부에서 실행 중입니다. 클러스터 이름을 자동으로 감지합니다.")

		// 방법 1: ServiceAccount 토큰에서 클러스터 이름 추출
		clusterName, err := getEksClusterNameFromServiceAccountToken(cfg)
		if err == nil && clusterName != "" {
			fmt.Printf("ServiceAccount 토큰에서 클러스터 이름을 찾았습니다: %s\n", clusterName)
			return clusterName
		} else if err != nil {
			fmt.Printf("ServiceAccount 토큰 방식 실패: %v\n", err)
		}

		// 방법 2: 환경 변수에서 가져오기
		if clusterName := os.Getenv("CLUSTER_NAME"); clusterName != "" {
			fmt.Printf("환경 변수에서 클러스터 이름을 찾았습니다: %s\n", clusterName)
			return clusterName
		}

		// 방법 3: ConfigMap에서 클러스터 정보 읽기
		clientset, err := kubernetes.NewForConfig(&kubeconfig)
		if err == nil {
			// kube-system 네임스페이스의 aws-auth ConfigMap 확인
			configMap, err := clientset.CoreV1().ConfigMaps("kube-system").Get(context.Background(), "aws-auth", metav1.GetOptions{})
			if err == nil && configMap != nil {
				if clusterName, ok := configMap.Data["cluster-name"]; ok {
					fmt.Printf("aws-auth ConfigMap에서 클러스터 이름을 찾았습니다: %s\n", clusterName)
					return clusterName
				}
			}

			// 또는 EKS 클러스터의 경우 아래와 같은 ConfigMap에서도 정보를 찾을 수 있음
			configMap, err = clientset.CoreV1().ConfigMaps("kube-system").Get(context.Background(), "cluster-info", metav1.GetOptions{})
			if err == nil && configMap != nil {
				if data, ok := configMap.Data["cluster.name"]; ok {
					fmt.Printf("cluster-info ConfigMap에서 클러스터 이름을 찾았습니다: %s\n", data)
					return data
				}
			}
		}

		// 방법 4: 노드 이름에서 추출 시도
		nodeName := os.Getenv("NODE_NAME")
		if strings.Contains(nodeName, "compute.internal") {
			parts := strings.Split(nodeName, ".")
			if len(parts) > 1 {
				clusterName := "eks-cluster-in-" + parts[1]
				fmt.Printf("노드 이름에서 클러스터 이름을 추정했습니다: %s\n", clusterName)
				return clusterName
			}
		}

		fmt.Println("클러스터 이름을 찾을 수 없어 기본값을 사용합니다: in-cluster")
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

// getEksClusterNameFromServiceAccountToken 함수 수정 - AWS 설정을 인자로 받도록 변경
func getEksClusterNameFromServiceAccountToken(cfg aws.Config) (string, error) {
	// 서비스 어카운트 토큰 파일 경로
	tokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"

	// 토큰 파일이 존재하는지 확인
	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		return "", fmt.Errorf("서비스 어카운트 토큰 파일을 찾을 수 없습니다: %v", err)
	}

	// 토큰 파일 읽기
	tokenBytes, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", fmt.Errorf("토큰 파일을 읽을 수 없습니다: %v", err)
	}

	// JWT 토큰은 헤더.페이로드.서명 형식
	tokenParts := strings.Split(string(tokenBytes), ".")
	if len(tokenParts) != 3 {
		return "", fmt.Errorf("유효하지 않은 JWT 토큰 형식입니다")
	}

	// 페이로드 부분만 추출 (Base64 인코딩된 상태)
	payloadBase64 := tokenParts[1]

	// 패딩 추가 (JWT에서는 패딩이 생략됨)
	if len(payloadBase64)%4 != 0 {
		payloadBase64 += strings.Repeat("=", 4-len(payloadBase64)%4)
	}

	// Base64 디코딩
	payloadBytes, err := base64.StdEncoding.DecodeString(payloadBase64)
	if err != nil {
		// URL-safe Base64 방식으로 재시도
		payloadBase64 = strings.ReplaceAll(payloadBase64, "-", "+")
		payloadBase64 = strings.ReplaceAll(payloadBase64, "_", "/")
		payloadBytes, err = base64.StdEncoding.DecodeString(payloadBase64)
		if err != nil {
			return "", fmt.Errorf("토큰 페이로드를 디코딩할 수 없습니다: %v", err)
		}
	}

	// JSON 파싱
	var payload TokenPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return "", fmt.Errorf("토큰 페이로드를 파싱할 수 없습니다: %v", err)
	}

	// iss 값 확인
	if payload.Iss == "" {
		return "", fmt.Errorf("토큰에 발급자(iss) 정보가 없습니다")
	}

	// iss 값에서 클러스터 ID 추출 (https://oidc.eks.{region}.amazonaws.com/id/{cluster-id} 형식)
	parts := strings.Split(payload.Iss, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("토큰의 발급자(iss) 형식이 예상과 다릅니다: %s", payload.Iss)
	}

	clusterId := parts[len(parts)-1]

	// 리전 추출 (oidc.eks.{region}.amazonaws.com 형식)
	issuerParts := strings.Split(payload.Iss, ".")
	if len(issuerParts) < 4 {
		return "", fmt.Errorf("토큰의 발급자(iss)에서 리전을 추출할 수 없습니다: %s", payload.Iss)
	}

	region := issuerParts[2]

	// 환경 변수에 리전 설정 (AWS CLI가 사용할 수 있게)
	os.Setenv("AWS_REGION", region)
	fmt.Printf("토큰에서 리전 정보를 찾아 설정했습니다: %s\n", region)

	// AWS SDK를 사용하여 클러스터 ID로 클러스터 이름 조회
	return getClusterNameByID(clusterId, region, cfg)
}

// getClusterNameByID 함수 수정 - AWS SDK를 사용하도록 변경
func getClusterNameByID(clusterId, region string, cfg aws.Config) (string, error) {
	// 리전 설정이 있으면 해당 리전으로 설정 업데이트
	if region != "" {
		cfg.Region = region
	}

	// EKS 클라이언트 생성
	eksClient := eks.NewFromConfig(cfg)

	// EKS 클러스터 목록 가져오기
	listClustersOutput, err := eksClient.ListClusters(context.TODO(), &eks.ListClustersInput{})
	if err != nil {
		return "", fmt.Errorf("EKS 클러스터 목록을 조회하는 중 오류 발생: %v", err)
	}

	// 각 클러스터에 대해 정보 조회하여 클러스터 ID 비교
	for _, clusterName := range listClustersOutput.Clusters {
		// 클러스터 상세 정보 가져오기
		describeOutput, err := eksClient.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
			Name: aws.String(clusterName),
		})

		if err != nil {
			fmt.Printf("클러스터 %s 정보 조회 중 오류 발생: %v\n", clusterName, err)
			continue
		}

		// OIDC 발급자 URL에 클러스터 ID가 포함되어 있는지 확인
		if describeOutput.Cluster != nil &&
			describeOutput.Cluster.Identity != nil &&
			describeOutput.Cluster.Identity.Oidc != nil &&
			describeOutput.Cluster.Identity.Oidc.Issuer != nil {

			issuer := *describeOutput.Cluster.Identity.Oidc.Issuer

			if strings.Contains(issuer, clusterId) {
				fmt.Printf("클러스터 ID %s에 해당하는 클러스터를 찾았습니다: %s\n", clusterId, clusterName)
				return clusterName, nil
			}
		}
	}

	return "", fmt.Errorf("클러스터 ID %s에 해당하는 EKS 클러스터를 찾을 수 없습니다", clusterId)
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

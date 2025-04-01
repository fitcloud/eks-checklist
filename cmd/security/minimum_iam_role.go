package security

import (
	"context"
	"fmt"
	"strings"

	"eks-checklist/cmd/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 데이터 플레인 노드에 허용된 IAM 정책 목록 (이 외의 정책은 비허용으로 간주)
var allowedPolicies = map[string]bool{
	"AmazonEC2ContainerRegistryReadOnly": true,
	"AmazonEKS_CNI_Policy":               true,
	"AmazonEKSWorkerNodePolicy":          true,
}

// GetNodeIPs는 모든 노드에서 제공된 IP 주소(provided-node-ip 어노테이션)를 수집
func GetNodeIPs(client kubernetes.Interface) ([]string, error) {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var nodeIPs []string
	for _, node := range nodes.Items {
		// 노드 어노테이션 중 provided-node-ip 값을 수집
		if ip, ok := node.Annotations["alpha.kubernetes.io/provided-node-ip"]; ok {
			nodeIPs = append(nodeIPs, ip)
		}
	}

	return nodeIPs, nil
}

// GetIAMRoleForNode는 주어진 노드 IP에 연결된 EC2 인스턴스의 IAM 역할 이름을 반환
func GetIAMRoleForNode(nodeIP string) (string, error) {
	// AWS 세션 초기화 (공유 구성 사용)
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := ec2.New(sess)

	// EC2 인스턴스 필터링: private IP로 인스턴스를 조회
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-ip-address"),
				Values: []*string{aws.String(nodeIP)},
			},
		},
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
		return "", err
	}

	// 조회된 인스턴스가 없는 경우 에러 반환
	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return "", fmt.Errorf("no instance found for IP %s", nodeIP)
	}

	instance := result.Reservations[0].Instances[0]

	// IAM 인스턴스 프로파일이 없는 경우
	if instance.IamInstanceProfile == nil || instance.IamInstanceProfile.Arn == nil {
		return "", fmt.Errorf("no IAM role associated with instance %s", *instance.InstanceId)
	}

	// IAM 인스턴스 프로파일 이름 추출 (ARN에서 마지막 부분)
	profileArn := *instance.IamInstanceProfile.Arn
	profileName := profileArn[strings.LastIndex(profileArn, "/")+1:]

	// IAM 서비스 클라이언트 생성
	iamSvc := iam.New(sess)

	// 인스턴스 프로파일에서 역할 정보를 조회
	profileInput := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
	}
	profileOutput, err := iamSvc.GetInstanceProfile(profileInput)
	if err != nil {
		return "", fmt.Errorf("failed to get IAM instance profile details: %v", err)
	}

	// 프로파일에 역할이 없는 경우 에러
	if len(profileOutput.InstanceProfile.Roles) == 0 {
		return "", fmt.Errorf("no IAM role found in instance profile %s", profileName)
	}

	// 역할 이름 반환
	return *profileOutput.InstanceProfile.Roles[0].RoleName, nil
}

// GetAttachedPolicies는 지정된 IAM 역할에 연결된 정책 이름 목록을 반환
func GetAttachedPolicies(roleName string) ([]string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := iam.New(sess)

	// 역할에 연결된 정책 나열
	input := &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	}

	result, err := svc.ListAttachedRolePolicies(input)
	if err != nil {
		return nil, err
	}

	var policies []string
	for _, policy := range result.AttachedPolicies {
		policies = append(policies, *policy.PolicyName)
	}

	return policies, nil
}

// CheckNodeIAMRoles는 모든 노드의 IAM 역할에 허용되지 않은 정책이 있는지 확인
func CheckNodeIAMRoles(client kubernetes.Interface) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "데이터 플레인 노드에 필수로 필요한 IAM 권한만 부여",
		Manual:     false,
		Passed:     true, // 기본적으로 통과 상태로 설정, 문제 발생 시 false로 변경
		FailureMsg: "일부 노드에서 허용되지 않은 IAM 정책이 발견되었습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image", // 문제가 있을 경우 참고할 Runbook 링크
	}

	// 노드 IP 목록 가져오기
	nodeIPs, err := GetNodeIPs(client)
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	// 각 노드에 대해 IAM 역할 및 정책 확인
	for _, ip := range nodeIPs {
		roleName, err := GetIAMRoleForNode(ip)
		if err != nil {
			result.Passed = false
			result.Resources = append(result.Resources, fmt.Sprintf("Node IP: %s, Error: %v", ip, err))
			continue
		}

		policies, err := GetAttachedPolicies(roleName)
		if err != nil {
			result.Passed = false
			result.Resources = append(result.Resources, fmt.Sprintf("Role: %s, Error: %v", roleName, err))
			continue
		}

		// 허용되지 않은 정책이 포함되어 있는지 검사
		for _, policy := range policies {
			if !allowedPolicies[policy] {
				result.Passed = false
				result.Resources = append(result.Resources, fmt.Sprintf("Role: %s, Unauthorized Policy: %s", roleName, policy))
			}
		}
	}

	return result
}

package security

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var allowedPolicies = map[string]bool{
	"AmazonEC2ContainerRegistryReadOnly": true,
	"AmazonEKS_CNI_Policy":               true,
	"AmazonEKSWorkerNodePolicy":          true,
}

const (
	Red    = "\033[31m" // 빨간색
	Green  = "\033[32m" // 초록색
	Yellow = "\033[33m" // 노란색
	Reset  = "\033[0m"  // 기본 색상으로 리셋
)

// GetNodeIPs retrieves the provided-node-ip annotations from all nodes.
func GetNodeIPs(client kubernetes.Interface) []string {
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to list nodes: %v", err)
	}

	var nodeIPs []string
	for _, node := range nodes.Items {
		if ip, ok := node.Annotations["alpha.kubernetes.io/provided-node-ip"]; ok {
			nodeIPs = append(nodeIPs, ip)
		}
	}

	return nodeIPs
}

// GetIAMRoleForNode retrieves the IAM role associated with a given node IP.
func GetIAMRoleForNode(nodeIP string) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := ec2.New(sess)

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

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return "", fmt.Errorf("no instance found for IP %s", nodeIP)
	}

	instance := result.Reservations[0].Instances[0]
	if instance.IamInstanceProfile == nil || instance.IamInstanceProfile.Arn == nil {
		return "", fmt.Errorf("no IAM role associated with instance %s", *instance.InstanceId)
	}

	profileArn := *instance.IamInstanceProfile.Arn
	profileName := profileArn[strings.LastIndex(profileArn, "/")+1:]

	iamSvc := iam.New(sess)
	profileInput := &iam.GetInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
	}

	profileOutput, err := iamSvc.GetInstanceProfile(profileInput)
	if err != nil {
		return "", fmt.Errorf("failed to get IAM instance profile details: %v", err)
	}

	if len(profileOutput.InstanceProfile.Roles) == 0 {
		return "", fmt.Errorf("no IAM role found in instance profile %s", profileName)
	}

	return *profileOutput.InstanceProfile.Roles[0].RoleName, nil
}

// GetAttachedPolicies fetches the attached IAM policies for a given role.
func GetAttachedPolicies(roleName string) ([]string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := iam.New(sess)

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

// CheckNodeIAMRoles fetches IAM roles for all nodes and verifies their attached policies.
func CheckNodeIAMRoles(client kubernetes.Interface) bool {
	nodeIPs := GetNodeIPs(client)
	for _, ip := range nodeIPs {
		//		fmt.Printf("🔍 Checking IAM role for node with IP: %s\n", ip)
		roleName, err := GetIAMRoleForNode(ip)
		if err != nil {
			fmt.Printf("⚠️  Failed to get IAM role for node IP %s: %v\n", ip, err)
			return false
		}

		//		fmt.Printf("ℹ️  Extracted IAM Role Name: %s\n", roleName)
		policies, err := GetAttachedPolicies(roleName)
		if err != nil {
			fmt.Printf("⚠️  Failed to get IAM policies for role %s: %v\n", roleName, err)
			return false
		}

		for _, policy := range policies {
			if !allowedPolicies[policy] {
				fmt.Printf(Red+"✖ FAIL: Unauthorized policy detected: %s on role %s\n"+Reset, policy, roleName)
				return false
			}
		}
	}

	fmt.Println(Green + "✔ PASS: All nodes are using only allowed IAM policies" + Reset)
	return true
}

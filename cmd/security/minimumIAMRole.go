package security

import (
	"context"
	"fmt"
	"log"
	"strings"

	//디펜던시에 추가 필요
	// go get -u github.com/aws/aws-sdk-go/...
	//
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
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-2"), // 필요에 따라 변경하세요.
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

	return *instance.IamInstanceProfile.Arn, nil
}

// GetAttachedPolicies fetches the attached IAM policies for a given role.
func GetAttachedPolicies(roleName string) ([]string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
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
		roleArn, err := GetIAMRoleForNode(ip)
		if err != nil {
			fmt.Printf("⚠️  Failed to get IAM role for node IP %s: %v\n", ip, err)
			return false
		}

		roleName := roleArn[strings.LastIndex(roleArn, "/")+1:]
		policies, err := GetAttachedPolicies(roleName)
		if err != nil {
			fmt.Printf("⚠️  Failed to get IAM policies for role %s: %v\n", roleName, err)
			return false
		}

		for _, policy := range policies {
			if !allowedPolicies[policy] {
				fmt.Printf("❌ Unauthorized policy detected: %s on role %s\n", policy, roleName)
				return false
			}
		}
	}

	fmt.Println("✅ All nodes have only allowed IAM policies.")
	return true
}

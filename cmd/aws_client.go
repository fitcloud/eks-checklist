package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

// AWSClient는 AWS 서비스와 상호작용하기 위한 인터페이스입니다.
// 이 인터페이스를 사용하면 실제 AWS API 호출을 모킹하여 테스트하기 쉬워집니다.
type AWSClient interface {
	// DescribeEKSCluster는 EKS 클러스터 정보를 반환합니다.
	DescribeEKSCluster(ctx context.Context, clusterName string) (*eks.DescribeClusterOutput, error)
	
	// ListNodegroups는 EKS 클러스터의 노드 그룹 목록을 반환합니다.
	ListNodegroups(ctx context.Context, clusterName string) ([]string, error)
	
	// DescribeNodegroup은 EKS 노드 그룹 정보를 반환합니다.
	DescribeNodegroup(ctx context.Context, clusterName, nodegroupName string) (*eks.DescribeNodegroupOutput, error)
}

// DefaultAWSClient는 실제 AWS SDK를 사용하는 AWSClient 구현체입니다.
type DefaultAWSClient struct {
	config aws.Config
}

// NewAWSClient는 새로운 DefaultAWSClient 인스턴스를 생성합니다.
func NewAWSClient(config aws.Config) *DefaultAWSClient {
	return &DefaultAWSClient{
		config: config,
	}
}

// DescribeEKSCluster는 EKS 클러스터 정보를 반환합니다.
func (c *DefaultAWSClient) DescribeEKSCluster(ctx context.Context, clusterName string) (*eks.DescribeClusterOutput, error) {
	client := eks.NewFromConfig(c.config)
	return client.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})
}

// ListNodegroups는 EKS 클러스터의 노드 그룹 목록을 반환합니다.
func (c *DefaultAWSClient) ListNodegroups(ctx context.Context, clusterName string) ([]string, error) {
	client := eks.NewFromConfig(c.config)
	result, err := client.ListNodegroups(ctx, &eks.ListNodegroupsInput{
		ClusterName: aws.String(clusterName),
	})
	if err != nil {
		return nil, err
	}
	return result.Nodegroups, nil
}

// DescribeNodegroup은 EKS 노드 그룹 정보를 반환합니다.
func (c *DefaultAWSClient) DescribeNodegroup(ctx context.Context, clusterName, nodegroupName string) (*eks.DescribeNodegroupOutput, error) {
	client := eks.NewFromConfig(c.config)
	return client.DescribeNodegroup(ctx, &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(clusterName),
		NodegroupName: aws.String(nodegroupName),
	})
}

// MockAWSClient는 테스트를 위한 AWSClient 모의 구현체입니다.
type MockAWSClient struct {
	MockDescribeEKSCluster func(ctx context.Context, clusterName string) (*eks.DescribeClusterOutput, error)
	MockListNodegroups     func(ctx context.Context, clusterName string) ([]string, error)
	MockDescribeNodegroup  func(ctx context.Context, clusterName, nodegroupName string) (*eks.DescribeNodegroupOutput, error)
}

// DescribeEKSCluster는 모의 EKS 클러스터 정보를 반환합니다.
func (m *MockAWSClient) DescribeEKSCluster(ctx context.Context, clusterName string) (*eks.DescribeClusterOutput, error) {
	if m.MockDescribeEKSCluster != nil {
		return m.MockDescribeEKSCluster(ctx, clusterName)
	}
	// 기본 모의 응답
	return &eks.DescribeClusterOutput{
		Cluster: &types.Cluster{
			Name: aws.String(clusterName),
		},
	}, nil
}

// ListNodegroups는 모의 EKS 클러스터의 노드 그룹 목록을 반환합니다.
func (m *MockAWSClient) ListNodegroups(ctx context.Context, clusterName string) ([]string, error) {
	if m.MockListNodegroups != nil {
		return m.MockListNodegroups(ctx, clusterName)
	}
	// 기본 모의 응답
	return []string{"mock-nodegroup"}, nil
}

// DescribeNodegroup은 모의 EKS 노드 그룹 정보를 반환합니다.
func (m *MockAWSClient) DescribeNodegroup(ctx context.Context, clusterName, nodegroupName string) (*eks.DescribeNodegroupOutput, error) {
	if m.MockDescribeNodegroup != nil {
		return m.MockDescribeNodegroup(ctx, clusterName, nodegroupName)
	}
	// 기본 모의 응답
	return &eks.DescribeNodegroupOutput{
		Nodegroup: &types.Nodegroup{
			NodegroupName: aws.String(nodegroupName),
			ClusterName:   aws.String(clusterName),
		},
	}, nil
} 
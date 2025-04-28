package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

type EksCluster struct {
	Cluster *types.Cluster
}

func Describe(clusterName string, cfg aws.Config) EksCluster {
	eksClient := eks.NewFromConfig(cfg)
	output, err := eksClient.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
		Name: &clusterName,
	})

	if err != nil {
		panic(err.Error())
	}

	eksCluster := EksCluster{Cluster: output.Cluster}

	return eksCluster
}

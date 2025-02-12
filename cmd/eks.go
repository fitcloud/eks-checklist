package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

type EksCluster struct {
	Cluster *types.Cluster
}

func Describe(clusterName string) EksCluster {
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err.Error())
	}

	eksClient := eks.NewFromConfig(awsConfig)
	output, err := eksClient.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
		Name: &clusterName,
	})

	if err != nil {
		panic(err.Error())
	}

	eksCluster := EksCluster{Cluster: output.Cluster}

	return eksCluster
}

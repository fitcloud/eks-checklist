package cmd

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	awsConfig aws.Config
	once      sync.Once
)

// AWS 설정을 로드하는 함수 (싱글톤)
func GetAWSConfig() aws.Config {
	once.Do(func() {
		var err error
		awsConfig, err = config.LoadDefaultConfig(context.TODO())
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}
	})
	return awsConfig
}

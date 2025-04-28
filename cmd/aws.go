package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	awsConfig aws.Config
	once      sync.Once
)

// AWS 설정을 로드하는 함수 (싱글톤)
func GetAWSConfig(AWS_PROFILE string) aws.Config {
	fmt.Printf("AWS_PROFILE: %s\n", AWS_PROFILE)
	if AWS_PROFILE != "" {
		once.Do(func() {
			var err error
			awsConfig, err = config.LoadDefaultConfig(
				context.TODO(),
				config.WithSharedConfigProfile(AWS_PROFILE),
			)
			if err != nil {
				log.Printf("AWS 프로필 '%s'로 설정을 로드할 수 없습니다: %v", AWS_PROFILE, err)
				// 실패 시 종료
				os.Exit(1)
			}
		})
		return awsConfig
	} else {
		once.Do(func() {
			var err error
			awsConfig, err = config.LoadDefaultConfig(context.TODO())
			if err != nil {
				log.Fatalf("unable to load SDK config, %v", err)
			}
		})
		return awsConfig
	}
}

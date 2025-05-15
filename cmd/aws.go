package cmd

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	awsConfig aws.Config
	once      sync.Once
	configErr error
)

// AWS 설정을 로드하는 함수 (싱글톤)
func GetAWSConfig(awsProfile string) (aws.Config, error) {
	fmt.Printf("AWS Profile: %s\n", awsProfile)
	
	// 이미 설정이 로드되었다면 저장된 결과 반환
	if awsConfig != (aws.Config{}) {
		return awsConfig, configErr
	}
	
	// 프로필이 지정된 경우
	if awsProfile != "" {
		once.Do(func() {
			awsConfig, configErr = config.LoadDefaultConfig(
				context.TODO(),
				config.WithSharedConfigProfile(awsProfile),
			)
			if configErr != nil {
				configErr = fmt.Errorf("AWS 프로필 '%s'로 설정을 로드할 수 없습니다: %w", awsProfile, configErr)
			}
		})
	} else {
		// 기본 프로필 사용
		once.Do(func() {
			awsConfig, configErr = config.LoadDefaultConfig(context.TODO())
			if configErr != nil {
				configErr = fmt.Errorf("AWS 설정을 로드할 수 없습니다: %w", configErr)
			}
		})
	}
	
	return awsConfig, configErr
}

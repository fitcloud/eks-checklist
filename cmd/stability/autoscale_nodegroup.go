// 기존 코드

// package stability

// import (
// 	"context"
// 	"log"
// 	"os"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/autoscaling"
// 	"github.com/aws/aws-sdk-go/service/eks"
// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// )

// // CheckAutoScaledManagedNodeGroup - ASG 기반 관리형 노드 그룹 자동 확장 여부 확인
// func CheckAutoScaledManagedNodeGroup(client kubernetes.Interface, clusterName string) bool {

// 	awsProfile := os.Getenv("AWS_PROFILE")
// 	region := os.Getenv("AWS_REGION")

// 	// AWS 세션 생성
// 	sess, err := session.NewSessionWithOptions(session.Options{
// 		Profile:           awsProfile,
// 		SharedConfigState: session.SharedConfigEnable,
// 		Config: aws.Config{
// 			Region: aws.String(region),
// 		},
// 	})
// 	if err != nil {
// 		log.Println("AWS 세션 생성 실패:", err)
// 		return false
// 	}

// 	eksClient := eks.New(sess)
// 	asgClient := autoscaling.New(sess)

// 	// 노드의 라벨을 확인하여 관리형 노드 그룹을 찾아서 map에 저장
// 	nodes, err := client.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
// 	if err != nil {
// 		log.Println("노드 목록 조회 실패:", err)
// 		return false
// 	}

// 	managedNodeGroups := make(map[string]bool)
// 	for _, node := range nodes.Items {
// 		if nodeGroup, ok := node.Labels["eks.amazonaws.com/nodegroup"]; ok {
// 			managedNodeGroups[nodeGroup] = true
// 		}
// 	}

// 	if len(managedNodeGroups) == 0 {
// 		log.Println("관리형 노드 그룹을 찾을 수 없음")
// 		return false
// 	}

// 	// 노드그룹들의 ASG 정보 확인
// 	for nodeGroup := range managedNodeGroups {
// 		ng, err := eksClient.DescribeNodegroup(&eks.DescribeNodegroupInput{
// 			ClusterName:   aws.String(clusterName),
// 			NodegroupName: aws.String(nodeGroup),
// 		})
// 		if err != nil {
// 			log.Println("노드 그룹 정보 조회 실패:", err)
// 			continue
// 		}

// 		if len(ng.Nodegroup.Resources.AutoScalingGroups) == 0 {
// 			log.Println("ASG에 의해 관리되지 않는 노드 그룹:", nodeGroup)
// 			continue
// 		}

// 		asgName := ng.Nodegroup.Resources.AutoScalingGroups[0].Name

// 		// ASG Scaling 설정 확인
// 		asg, err := asgClient.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
// 			AutoScalingGroupNames: []*string{asgName},
// 		})
// 		if err != nil || len(asg.AutoScalingGroups) == 0 {
// 			log.Println("ASG 정보 조회 실패:", err)
// 			continue
// 		}

// 		// 설정을 확인해서 minSize < maxSize 인 경우만 자동 확장으로 판단 왜냐하면 같거나 큰 경우는 자동이 아니니까
// 		asgConfig := asg.AutoScalingGroups[0]
// 		if asgConfig.MinSize != nil && asgConfig.MaxSize != nil && *asgConfig.MinSize < *asgConfig.MaxSize {
// 			log.Println("관리형 노드 그룹", nodeGroup, "은 ASG 기반으로 자동 확장됨")
// 			return true
// 		} else {
// 			log.Println("관리형 노드 그룹", nodeGroup, "은 ASG Scaling 설정이 없음")
// 		}
// 	}

// 	return false
// }

// 변경 후 코드
package stability

import (
	"context"
	"fmt"
	"os"

	"eks-checklist/cmd/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/eks"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckAutoScaledManagedNodeGroup - ASG 기반 관리형 노드 그룹 자동 확장 여부 확인
func CheckAutoScaledManagedNodeGroup(client kubernetes.Interface, clusterName string) common.CheckResult {
	result := common.CheckResult{
		CheckName:  "오토스케일링 그룹 기반 관리형 노드 그룹 생성",
		Manual:     false,
		Passed:     true,
		FailureMsg: "일부 관리형 노드 그룹이 ASG를 통한 자동 확장 구성이 되어 있지 않습니다.",
		Runbook:    "https://your.runbook.url/latest-tag-image",
	}

	result.Passed = false // 기본은 실패

	awsProfile := os.Getenv("AWS_PROFILE")
	region := os.Getenv("AWS_REGION")

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile:           awsProfile,
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(region),
		},
	})
	if err != nil {
		result.Passed = false
		result.FailureMsg = result.CheckName + " 검사 실패 : " + err.Error()
		return result
	}

	eksClient := eks.New(sess)
	asgClient := autoscaling.New(sess)

	nodes, err := client.CoreV1().Nodes().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		result.FailureMsg = fmt.Sprintf("노드 목록 조회 실패: %v", err)
		return result
	}

	managedNodeGroups := make(map[string]bool)
	for _, node := range nodes.Items {
		if nodeGroup, ok := node.Labels["eks.amazonaws.com/nodegroup"]; ok {
			managedNodeGroups[nodeGroup] = true
		}
	}

	if len(managedNodeGroups) == 0 {
		result.FailureMsg = "관리형 노드 그룹을 찾을 수 없습니다."
		return result
	}

	var (
		autoScaledCount  int
		totalNodeGroups  = len(managedNodeGroups)
		nonAutoScaledNgs []string
	)

	for nodeGroup := range managedNodeGroups {
		ng, err := eksClient.DescribeNodegroup(&eks.DescribeNodegroupInput{
			ClusterName:   aws.String(clusterName),
			NodegroupName: aws.String(nodeGroup),
		})
		if err != nil {
			nonAutoScaledNgs = append(nonAutoScaledNgs, nodeGroup+" (조회 실패)")
			continue
		}

		if len(ng.Nodegroup.Resources.AutoScalingGroups) == 0 {
			nonAutoScaledNgs = append(nonAutoScaledNgs, nodeGroup+" (ASG 없음)")
			continue
		}

		asgName := ng.Nodegroup.Resources.AutoScalingGroups[0].Name
		asg, err := asgClient.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{
			AutoScalingGroupNames: []*string{asgName},
		})
		if err != nil || len(asg.AutoScalingGroups) == 0 {
			nonAutoScaledNgs = append(nonAutoScaledNgs, nodeGroup+" (ASG 조회 실패)")
			continue
		}

		asgConf := asg.AutoScalingGroups[0]
		if asgConf.MinSize != nil && asgConf.MaxSize != nil && *asgConf.MinSize < *asgConf.MaxSize {
			autoScaledCount++
			result.Resources = append(result.Resources,
				fmt.Sprintf("Nodegroup: %s | ASG: %s (minSize: %d, maxSize: %d)",
					nodeGroup, *asgName, *asgConf.MinSize, *asgConf.MaxSize))
		} else {
			nonAutoScaledNgs = append(nonAutoScaledNgs, nodeGroup+" (minSize ≥ maxSize)")
		}
	}

	switch {
	case autoScaledCount == 0:
		result.Passed = false
	case autoScaledCount < totalNodeGroups:
		result.Passed = false
	default:
		result.Passed = true
	}

	return result
}

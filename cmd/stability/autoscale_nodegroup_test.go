package stability_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"eks-checklist/cmd/stability"
	"eks-checklist/cmd/testutils"

	"bou.ke/monkey"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/eks"

	// AWS SDK v1에서는 eks 타입은 바로 eks.AutoScalingGroup, eks.Nodegroup 등 사용
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckAutoScaledManagedNodeGroup(t *testing.T) {
	// YAML 파일 "autoscale_nodegroup.yaml"에서 테스트 케이스 로드
	testCases := testutils.LoadTestCases(t, "autoscale_nodegroup.yaml")
	for _, tc := range testCases {
		testName := tc["name"].(string)
		var expectedPass bool
		if v, ok := tc["expect_pass"]; ok && v != nil {
			expectedPass = v.(bool)
		} else {
			t.Fatalf("Test case '%s' missing 'expect_pass' field or it is nil", testName)
		}

		// cluster 정보 구성
		clusterMap, ok := tc["cluster"].(map[string]interface{})
		if !ok {
			t.Fatalf("Test case '%s' missing 'cluster' field", testName)
		}
		clusterName, ok := clusterMap["clusterName"].(string)
		if !ok {
			t.Fatalf("Test case '%s' missing 'clusterName' in cluster", testName)
		}
		resVpcCfg, ok := clusterMap["resourcesVpcConfig"].(map[string]interface{})
		if !ok {
			t.Fatalf("Test case '%s' missing 'resourcesVpcConfig' in cluster", testName)
		}
		subnetIdsRaw, ok := resVpcCfg["subnetIds"].([]interface{})
		if !ok {
			t.Fatalf("Test case '%s' missing 'subnetIds' in resourcesVpcConfig", testName)
		}
		var subnetIds []string
		for _, id := range subnetIdsRaw {
			s, ok := id.(string)
			if !ok {
				t.Fatalf("Test case '%s': subnetId is not a string", testName)
			}
			subnetIds = append(subnetIds, s)
		}
		// []string -> []*string 변환
		var ptrSubnetIds []*string
		for _, s := range subnetIds {
			ptrSubnetIds = append(ptrSubnetIds, aws.String(s))
		}

		// 노드 정보 구성 (노드의 Label에 "eks.amazonaws.com/nodegroup" 사용)
		nodesRaw, ok := tc["nodes"].([]interface{})
		if !ok {
			t.Fatalf("Test case '%s' missing 'nodes' field", testName)
		}

		// eks_nodegroups: mapping: nodegroupName -> { autoscaling_groups: [ { name: string } ] }
		eksNGRaw, ok := tc["eks_nodegroups"].(map[string]interface{})
		if !ok {
			t.Fatalf("Test case '%s' missing 'eks_nodegroups' field", testName)
		}
		eksNG := make(map[string]eks.Nodegroup)
		for ngName, raw := range eksNGRaw {
			m, ok := raw.(map[string]interface{})
			if !ok {
				t.Fatalf("Test case '%s': eks_nodegroups entry for %s is not a map", testName, ngName)
			}
			agRaw, ok := m["autoscaling_groups"].([]interface{})
			if !ok {
				t.Fatalf("Test case '%s': eks_nodegroups entry for %s missing 'autoscaling_groups'", testName, ngName)
			}
			var agSlice []eks.AutoScalingGroup
			for _, ag := range agRaw {
				agMap, ok := ag.(map[string]interface{})
				if !ok {
					t.Fatalf("Test case '%s': autoscaling_groups element is not a map", testName)
				}
				agName, ok := agMap["name"].(string)
				if !ok {
					t.Fatalf("Test case '%s': autoscaling_groups element missing 'name'", testName)
				}
				agSlice = append(agSlice, eks.AutoScalingGroup{
					Name: aws.String(agName),
				})
			}
			// []eks.AutoScalingGroup -> []*eks.AutoScalingGroup 변환
			var ptrAgs []*eks.AutoScalingGroup
			for i := range agSlice {
				ptrAgs = append(ptrAgs, &agSlice[i])
			}
			eksNG[ngName] = eks.Nodegroup{
				Resources: &eks.NodegroupResources{
					AutoScalingGroups: ptrAgs,
				},
			}
		}

		// asg: mapping asgName -> { minSize: number, maxSize: number }
		asgRaw, ok := tc["asg"].(map[string]interface{})
		if !ok {
			t.Fatalf("Test case '%s' missing 'asg' field", testName)
		}
		asgInfo := make(map[string]struct {
			MinSize int
			MaxSize int
		})
		for asgName, raw := range asgRaw {
			m, ok := raw.(map[string]interface{})
			if !ok {
				t.Fatalf("Test case '%s': asg entry for %s is not a map", testName, asgName)
			}
			// minSize 처리: float64 또는 int 모두 지원
			var minSize float64
			switch v := m["minSize"].(type) {
			case float64:
				minSize = v
			case int:
				minSize = float64(v)
			default:
				t.Fatalf("Test case '%s': asg entry for %s missing or invalid 'minSize'", testName, asgName)
			}
			// maxSize 처리
			var maxSize float64
			switch v := m["maxSize"].(type) {
			case float64:
				maxSize = v
			case int:
				maxSize = float64(v)
			default:
				t.Fatalf("Test case '%s': asg entry for %s missing or invalid 'maxSize'", testName, asgName)
			}
			asgInfo[asgName] = struct {
				MinSize int
				MaxSize int
			}{MinSize: int(minSize), MaxSize: int(maxSize)}
		}

		t.Run(testName, func(t *testing.T) {
			client := fake.NewSimpleClientset()

			// 노드 생성 (노드의 Label에 "eks.amazonaws.com/nodegroup" 사용)
			for _, n := range nodesRaw {
				nMap, ok := n.(map[string]interface{})
				if !ok {
					t.Fatalf("Test case '%s': node is not a map", testName)
				}
				name, ok := nMap["name"].(string)
				if !ok {
					t.Fatalf("Test case '%s': node missing 'name'", testName)
				}
				ng, ok := nMap["nodegroup"].(string)
				if !ok {
					t.Fatalf("Test case '%s': node missing 'nodegroup'", testName)
				}
				ip, ok := nMap["provided_node_ip"].(string)
				if !ok {
					t.Fatalf("Test case '%s': node missing 'provided_node_ip'", testName)
				}
				node := &corev1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: name,
						Annotations: map[string]string{
							"alpha.kubernetes.io/provided-node-ip": ip,
						},
						Labels: map[string]string{
							"eks.amazonaws.com/nodegroup": ng,
						},
					},
				}
				_, err := client.CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})
				if err != nil {
					t.Fatalf("Test case '%s': failed to create node %s: %v", testName, name, err)
				}
			}

			// Patch DescribeNodegroup: eks.EKS의 DescribeNodegroup 메서드를 패치하여 가짜 응답 반환
			patch1 := monkey.PatchInstanceMethod(reflect.TypeOf(new(eks.EKS)), "DescribeNodegroup",
				func(e *eks.EKS, input *eks.DescribeNodegroupInput) (*eks.DescribeNodegroupOutput, error) {
					ngName := aws.StringValue(input.NodegroupName)
					fakeNG, exists := eksNG[ngName]
					if !exists {
						return nil, fmt.Errorf("no fake nodegroup info for %s", ngName)
					}
					return &eks.DescribeNodegroupOutput{
						Nodegroup: &fakeNG,
					}, nil
				})
			defer patch1.Unpatch()

			// Patch DescribeAutoScalingGroups: autoscaling.AutoScaling의 DescribeAutoScalingGroups 메서드를 패치하여 가짜 응답 반환
			patch2 := monkey.PatchInstanceMethod(reflect.TypeOf(new(autoscaling.AutoScaling)), "DescribeAutoScalingGroups",
				func(a *autoscaling.AutoScaling, input *autoscaling.DescribeAutoScalingGroupsInput) (*autoscaling.DescribeAutoScalingGroupsOutput, error) {
					if len(input.AutoScalingGroupNames) == 0 {
						return nil, fmt.Errorf("no asg names provided")
					}
					asgName := aws.StringValue(input.AutoScalingGroupNames[0])
					info, exists := asgInfo[asgName]
					if !exists {
						return nil, fmt.Errorf("no fake asg info for %s", asgName)
					}
					asgObj := autoscaling.Group{
						MinSize: aws.Int64(int64(info.MinSize)),
						MaxSize: aws.Int64(int64(info.MaxSize)),
					}
					return &autoscaling.DescribeAutoScalingGroupsOutput{
						AutoScalingGroups: []*autoscaling.Group{&asgObj},
					}, nil
				})
			defer patch2.Unpatch()

			result := stability.CheckAutoScaledManagedNodeGroup(client, clusterName)
			if result.Passed != expectedPass {
				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expectedPass, result.Passed)
			}
		})
	}
}

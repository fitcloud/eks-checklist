package security_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"eks-checklist/cmd/security"
// 	"eks-checklist/cmd/testutils"

// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes/fake"

// 	"bou.ke/monkey"
// )

// func TestCheckNodeIAMRoles_YAML(t *testing.T) {
// 	// YAML 파일 "check_node_iam_roles.yaml"에서 테스트 케이스 로드
// 	testCases := testutils.LoadTestCases(t, "check_node_iam_roles.yaml")
// 	for _, tc := range testCases {
// 		testName, ok := tc["name"].(string)
// 		if !ok {
// 			t.Fatalf("Test case missing 'name' field")
// 		}
// 		expected, ok := tc["expected"].(bool)
// 		if !ok {
// 			t.Fatalf("Test case '%s' missing 'expected' field or not bool", testName)
// 		}

// 		// YAML의 nodes 항목 읽기
// 		nodesRaw, ok := tc["nodes"].([]interface{})
// 		if !ok {
// 			t.Fatalf("Test case '%s' missing 'nodes' field", testName)
// 		}

// 		// YAML의 iam_roles 항목 읽기
// 		iamRolesRaw, ok := tc["iam_roles"].(map[string]interface{})
// 		if !ok {
// 			t.Fatalf("Test case '%s' missing 'iam_roles' field", testName)
// 		}
// 		// iam_roles: map[string] => { role: string, policies: []string }
// 		type IamRoleInfo struct {
// 			Role     string
// 			Policies []string
// 		}
// 		iamRoles := make(map[string]IamRoleInfo)
// 		for ip, raw := range iamRolesRaw {
// 			m, ok := raw.(map[string]interface{})
// 			if !ok {
// 				t.Fatalf("Test case '%s': iam_roles for ip %s is not a map", testName, ip)
// 			}
// 			role, ok := m["role"].(string)
// 			if !ok {
// 				t.Fatalf("Test case '%s': iam_roles for ip %s missing role", testName, ip)
// 			}
// 			policiesRaw, ok := m["policies"].([]interface{})
// 			if !ok {
// 				t.Fatalf("Test case '%s': iam_roles for ip %s missing policies", testName, ip)
// 			}
// 			var policies []string
// 			for _, p := range policiesRaw {
// 				pStr, ok := p.(string)
// 				if !ok {
// 					t.Fatalf("Test case '%s': policy value is not a string", testName)
// 				}
// 				policies = append(policies, pStr)
// 			}
// 			iamRoles[ip] = IamRoleInfo{
// 				Role:     role,
// 				Policies: policies,
// 			}
// 		}

// 		t.Run(testName, func(t *testing.T) {
// 			client := fake.NewSimpleClientset()

// 			// YAML에 정의된 Node 객체 생성
// 			// Node는 클러스터 범위 리소스이므로 네임스페이스는 설정하지 않습니다.
// 			for _, n := range nodesRaw {
// 				nodeMap, ok := n.(map[string]interface{})
// 				if !ok {
// 					t.Fatalf("Test case '%s': node is not a map", testName)
// 				}
// 				// YAML에 namespace 필드가 있더라도 무시합니다.
// 				name, ok := nodeMap["name"].(string)
// 				if !ok {
// 					t.Fatalf("Test case '%s': node missing 'name'", testName)
// 				}
// 				ip, ok := nodeMap["provided_node_ip"].(string)
// 				if !ok {
// 					t.Fatalf("Test case '%s': node missing 'provided_node_ip'", testName)
// 				}
// 				node := &corev1.Node{
// 					ObjectMeta: metav1.ObjectMeta{
// 						Name: name,
// 						Annotations: map[string]string{
// 							"alpha.kubernetes.io/provided-node-ip": ip,
// 						},
// 					},
// 				}
// 				_, err := client.CoreV1().Nodes().Create(context.TODO(), node, metav1.CreateOptions{})
// 				if err != nil {
// 					t.Fatalf("Test case '%s': failed to create node %s: %v", testName, name, err)
// 				}
// 			}

// 			// Patch GetIAMRoleForNode: node IP를 키로 iamRoles mapping에서 Role을 반환하도록 함.
// 			patch1 := monkey.Patch(security.GetIAMRoleForNode, func(nodeIP string) (string, error) {
// 				info, exists := iamRoles[nodeIP]
// 				if !exists {
// 					return "", fmt.Errorf("no iam role info for node ip %s", nodeIP)
// 				}
// 				return info.Role, nil
// 			})
// 			defer patch1.Unpatch()

// 			// Patch GetAttachedPolicies: role 이름에 해당하는 정책 목록을 반환하도록 함.
// 			patch2 := monkey.Patch(security.GetAttachedPolicies, func(roleName string) ([]string, error) {
// 				for _, info := range iamRoles {
// 					if info.Role == roleName {
// 						return info.Policies, nil
// 					}
// 				}
// 				return nil, fmt.Errorf("no policies for role %s", roleName)
// 			})
// 			defer patch2.Unpatch()

// 			// 함수 실행 및 반환값 비교
// 			result := security.CheckNodeIAMRoles(client)
// 			if result.Passed != expected {
// 				t.Errorf("Test '%s' failed: expected %v, got %v", testName, expected, result)
// 			}
// 		})
// 	}
// }

package security

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// EndpointSlicesCheck 함수는 클러스터 내 모든 Pod의 컨테이너에 대해
// readOnlyRootFilesystem=true 설정 여부를 점검합니다.
// 단, Windows 노드에서 실행 중이거나 kube-system 네임스페이스에 있는 경우는 제외합니다.

// readOnlyRootFilesystem 설정 점검 결과를 담는 구조체
type CheckResult struct {
	Namespace string
	Pod       string
	Container string
	Message   string
	Status    string // Passed, Failed, Skipped
}

func ReadnonlyFilesystemCheck(client kubernetes.Interface) bool {
	var results []CheckResult
	resultIsValid := true

	// 모든 네임스페이스의 Pod 리스트 조회
	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// 노드 OS 정보를 캐싱해두기 위한 맵 (중복 호출 방지)
	nodeOSCache := make(map[string]string)

	for _, pod := range pods.Items {
		// kube-system 네임스페이스는 검사 대상에서 제외
		if pod.Namespace == "kube-system" {
			for _, container := range pod.Spec.Containers {
				results = append(results, CheckResult{
					Namespace: pod.Namespace,
					Pod:       pod.Name,
					Container: container.Name,
					Message:   "Pod is in kube-system namespace, skipping check",
					Status:    "Skipped",
				})
			}
			continue
		}

		nodeName := pod.Spec.NodeName
		var nodeOS string

		// 캐시에 있으면 가져오고, 없으면 노드에서 조회
		if cached, ok := nodeOSCache[nodeName]; ok {
			nodeOS = cached
		} else {
			node, err := client.CoreV1().Nodes().Get(context.TODO(), nodeName, v1.GetOptions{})
			if err != nil {
				log.Printf("Failed to get node %s for pod %s/%s: %v", nodeName, pod.Namespace, pod.Name, err)
				nodeOS = "unknown"
			} else {
				// 노드의 OS 정보는 레이블에서 확인: kubernetes.io/os
				if osLabel, exists := node.Labels["kubernetes.io/os"]; exists {
					nodeOS = osLabel
				} else {
					nodeOS = "unknown"
				}
			}
			nodeOSCache[nodeName] = nodeOS
		}

		// Iterate containers
		for _, container := range pod.Spec.Containers {
			// Windows 노드에서 실행 중인 컨테이너는 검사 생략
			if nodeOS == "windows" {
				results = append(results, CheckResult{
					Namespace: pod.Namespace,
					Pod:       pod.Name,
					Container: container.Name,
					Message:   "Node OS is 'windows', skipping check",
					Status:    "Skipped",
				})
				continue
			}

			sc := container.SecurityContext
			if sc == nil || sc.ReadOnlyRootFilesystem == nil || !*sc.ReadOnlyRootFilesystem {
				results = append(results, CheckResult{
					Namespace: pod.Namespace,
					Pod:       pod.Name,
					Container: container.Name,
					Message:   "readOnlyRootFilesystem is not set to true",
					Status:    "Failed",
				})
				resultIsValid = false
			} else {
				results = append(results, CheckResult{
					Namespace: pod.Namespace,
					Pod:       pod.Name,
					Container: container.Name,
					Message:   "readOnlyRootFilesystem is set to true",
					Status:    "Passed",
				})
			}
		}
	}

	printResults(results)
	return resultIsValid
}

func printResults(results []CheckResult) {
	var failed []CheckResult
	for _, res := range results {
		if res.Status == "Failed" {
			failed = append(failed, res)
		}
	}

	if len(failed) == 0 {
		fmt.Println("✅ PASS: All pods use readOnlyRootFilesystem=true.")
	} else {
		fmt.Println("❌ FAIL: Some containers do not use readOnlyRootFilesystem=true.")
		for _, res := range failed {
			fmt.Printf("- Namespace: %s | Pod: %s | Container: %s\n", res.Namespace, res.Pod, res.Container)
		}
		fmt.Println("Runbook URL: https://your-runbook-url-here")
	}
}

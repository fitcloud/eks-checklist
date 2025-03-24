package security

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckResult represents the result of a readOnlyRootFilesystem check.
type CheckResult struct {
	Namespace string
	Pod       string
	Container string
	Message   string
	Status    string // Passed, Failed, Skipped
}

// EndpointSlicesCheck checks whether all containers in the cluster have readOnlyRootFilesystem=true,
// except for those running on Windows nodes.
func ReadnonlyFilesystemCheck(client kubernetes.Interface) {
	var results []CheckResult

	// List all pods in all namespaces
	pods, err := client.CoreV1().Pods("").List(context.TODO(), v1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Cache to store node OS info
	nodeOSCache := make(map[string]string)

	for _, pod := range pods.Items {
		nodeName := pod.Spec.NodeName
		var nodeOS string

		// Get node OS from cache or fetch from cluster
		if cached, ok := nodeOSCache[nodeName]; ok {
			nodeOS = cached
		} else {
			node, err := client.CoreV1().Nodes().Get(context.TODO(), nodeName, v1.GetOptions{})
			if err != nil {
				log.Printf("Failed to get node %s for pod %s/%s: %v", nodeName, pod.Namespace, pod.Name, err)
				nodeOS = "unknown"
			} else {
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

	PrintReadOnlyRootFSResults(results)
}

// PrintReadOnlyRootFSResults prints the results in a human-readable format.
func PrintReadOnlyRootFSResults(results []CheckResult) {
	var failed []CheckResult
	var skipped []CheckResult

	for _, res := range results {
		switch res.Status {
		case "Failed":
			failed = append(failed, res)
		case "Skipped":
			skipped = append(skipped, res)
		}
	}

	if len(failed) == 0 {
		fmt.Println(Green + "PASS: All pods use readOnlyRootFilesystem=true." + Reset)
	} else {
		fmt.Println(Red + "FAIL: Some containers do not use readOnlyRootFilesystem=true." + Reset)
		fmt.Println("Affected resources:")
		for _, res := range failed {
			fmt.Printf("- Namespace: %s | Pod: %s | Container: %s\n", res.Namespace, res.Pod, res.Container)
		}
		fmt.Println("Runbook URL: https://your-runbook-url-here")
	}

	// if len(skipped) > 0 {
	// 	fmt.Println()
	// 	fmt.Println("SKIPPED: Some containers were skipped because they run on Windows nodes.")
	// 	for _, res := range skipped {
	// 		fmt.Printf("- Namespace: %s | Pod: %s | Container: %s\n", res.Namespace, res.Pod, res.Container)
	// 	}
	// }
}

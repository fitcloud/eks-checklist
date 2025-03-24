package security

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CheckResult1 struct {
	Namespace        string
	ServiceAccount   string
	UsingIRSA        bool
	UsingPodIdentity bool
}

func CheckIRSAAndPodIdentity(clientset kubernetes.Interface) {
	// 로컬 구조체 정의
	type CheckResult struct {
		Namespace        string
		ServiceAccount   string
		UsingIRSA        bool
		UsingPodIdentity bool
	}

	// 내부 출력 함수 정의
	printResults := func(results []CheckResult) {
		var affected []string

		for _, res := range results {
			if !res.UsingIRSA && !res.UsingPodIdentity {
				affected = append(affected, fmt.Sprintf("- %s/%s", res.Namespace, res.ServiceAccount))
			}
		}

		if len(affected) == 0 {
			fmt.Println(Green + "PASS : All service accounts in this cluster are using either IRSA or EKS Pod Identity." + Reset)
		} else {
			fmt.Println(Red + "FAIL : Some service accounts are not configured with IRSA or EKS Pod Identity." + Reset)
			fmt.Println("Affected service accounts:")
			for _, sa := range affected {
				fmt.Println(sa)
			}
			fmt.Println("Runbook URL: https://your.runbook.url/irsa-or-pod-identity")
		}
	}

	// 체크 로직
	// 모든 SA가 아닌 주요 워크로드와 오픈소스들의 SA에 해당 부분을 검사해야하는데 기준을 아직 잡지 못하여 kube-system namespace만 제외 하여 검색
	saList, err := clientset.CoreV1().ServiceAccounts("").List(context.TODO(), v1.ListOptions{
		FieldSelector: "metadata.namespace!=kube-system", // kube-system 네임스페이스 제외
	})
	if err != nil {
		panic(err.Error())
	}

	var results []CheckResult
	// 각 ServiceAccount에 대해 IRSA 또는 Pod Identity 사용 여부 검사
	for _, sa := range saList.Items {
		annotations := sa.Annotations

		_, hasIRSA := annotations["eks.amazonaws.com/role-arn"]
		_, hasIdentity := annotations["eks.amazonaws.com/identity"]
		_, hasAudience := annotations["eks.amazonaws.com/audience"]

		results = append(results, CheckResult{
			Namespace:        sa.Namespace,
			ServiceAccount:   sa.Name,
			UsingIRSA:        hasIRSA,
			UsingPodIdentity: hasIdentity || hasAudience,
		})
	}

	printResults(results)
}

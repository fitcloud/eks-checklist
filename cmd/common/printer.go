package common

import "fmt"

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
)

var (
	PassedCount int
	FailedCount int
	ManualCount int
)

func PrintResult(r CheckResult) {
	// 필터 기준에 따라 이 결과를 출력할지 확인
	if !ShouldPrintResult(r.Passed, r.Manual) {
		return // 이 결과는 출력하지 않음
	}

	if r.Passed {
		PassedCount++
	} else if r.Manual {
		ManualCount++
	} else {
		FailedCount++
	}

	if r.Passed {
		fmt.Printf(Green+"✔ PASS | %s\n"+Reset, r.CheckName)
	} else {
		if r.Manual {
			fmt.Printf(Yellow+"⚠ MANUAL | %s\n"+Reset, r.CheckName)
		} else {
			fmt.Printf(Red+"✖ FAIL | %s\n"+Reset, r.CheckName)
		}
		fmt.Printf("  ├─ 🔸 이유 : %s\n", r.FailureMsg)
		if len(r.Resources) > 0 {
			fmt.Printf("  ├─ 🔸 영향받는 리소스:\n")
			for _, res := range r.Resources {
				fmt.Printf("  │   └─ %s\n", res)
			}
		}
		fmt.Printf("  └─ 🔗 Runbook: %s\n", r.Runbook)
	}
	fmt.Println()
}

func PrintSummary() {
	fmt.Println("\n===============[Checklist Summary]===============")
	fmt.Printf(Green+"✔ PASS: %d\n"+Reset, PassedCount)
	fmt.Printf(Red+"✖ FAIL: %d\n"+Reset, FailedCount)
	fmt.Printf(Yellow+"⚠ Manual: %d\n"+Reset, ManualCount)
	fmt.Println("===============[End of Summary]=================")
}

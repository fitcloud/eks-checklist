package common

import "fmt"

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
)

var (
	PassedCount     int
	FailedCount     int
	ManualCount     int
	CurrentCategory string
)

// SetCurrentCategory 현재 카테고리 설정
func SetCurrentCategory(category string) {
	CurrentCategory = category
}

// PrintCategoryHeader 카테고리 헤더 출력
func PrintCategoryHeader(category string) {
	SetCurrentCategory(category)

	if OutputFormat == "html" || OutputFormat == "pdf" {
		return
	}

	fmt.Printf("\n===============[%s]===============\n", category)
}

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

	// HTML 출력을 위한 결과 추가
	if OutputFormat == "html" || OutputFormat == "pdf" {
		AddResultForHTML(r, CurrentCategory)
		return
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
	if OutputFormat == "html" || OutputFormat == "pdf" {
		// HTML 보고서 저장
		htmlFilePath, err := SaveHTMLReport()
		if err != nil {
			fmt.Printf("HTML 보고서 생성 오류: %v\n", err)
			return
		}

		if OutputFormat == "html" {
			fmt.Printf("HTML 보고서가 %s에 저장되었습니다.\n", htmlFilePath)
			return // HTML 보고서 저장 후 종료
		}

		// PDF 변환이 필요한 경우
		if OutputFormat == "pdf" {
			pdfFilePath, err := ConvertHTMLToPDF(htmlFilePath)
			if err != nil {
				fmt.Printf("PDF 변환 오류: %v\n", err)
				return
			}
			fmt.Printf("PDF 보고서가 %s에 저장되었습니다.\n", pdfFilePath)
		}
		return
	}

	fmt.Println("\n===============[Checklist Summary]===============")
	fmt.Printf(Green+"✔ PASS: %d\n"+Reset, PassedCount)
	fmt.Printf(Red+"✖ FAIL: %d\n"+Reset, FailedCount)
	fmt.Printf(Yellow+"⚠ Manual: %d\n"+Reset, ManualCount)
	fmt.Println("===============[End of Summary]=================")
}

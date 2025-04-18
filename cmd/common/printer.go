package common

import (
	"fmt"
	"sort"
)

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

	// 정렬 모드 관련 변수들
	SortByStatus      bool              // 상태별 정렬 여부
	sortedResults     []CheckResult     // 정렬 모드에서 결과를 임시 저장
	sortedHtmlResults []CheckResultHTML // HTML 출력용 정렬된 결과
)

// SetSortMode 정렬 모드 설정
func SetSortMode(sortMode bool) {
	SortByStatus = sortMode
	if SortByStatus {
		// 정렬 모드가 활성화되면 결과 저장 컨테이너 초기화
		sortedResults = []CheckResult{}
		sortedHtmlResults = []CheckResultHTML{}
	}
}

// SetCurrentCategory 현재 카테고리 설정
func SetCurrentCategory(category string) {
	CurrentCategory = category
}

// PrintCategoryHeader 카테고리 헤더 출력
func PrintCategoryHeader(category string) {
	SetCurrentCategory(category)

	// 정렬 모드이거나 HTML/PDF 출력 모드인 경우 헤더를 출력하지 않음
	if SortByStatus || OutputFormat == "html" || OutputFormat == "pdf" {
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
		// 정렬 모드일 경우 결과를 바로 추가하지 않고 저장
		if SortByStatus {
			// 카테고리 정보를 결과에 저장
			r.Category = CurrentCategory
			sortedResults = append(sortedResults, r)
			return
		}
		AddResultForHTML(r, CurrentCategory)
		return
	}

	// 정렬 모드일 경우 결과를 바로 출력하지 않고 저장
	if SortByStatus {
		// 카테고리 정보를 결과에 저장
		r.Category = CurrentCategory
		sortedResults = append(sortedResults, r)
		return
	}

	// 일반 텍스트 출력
	printSingleResult(r)
}

// printSingleResult 단일 결과 출력
func printSingleResult(r CheckResult) {
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
		// 정렬 모드에서는 카테고리 정보도 출력
		if SortByStatus && r.Category != "" {
			fmt.Printf("      📂 카테고리: %s\n", r.Category)
		}
	}
	fmt.Println()
}

func PrintSummary() {
	// 정렬 모드이고 텍스트 출력인 경우 저장된 결과를 상태별로 출력
	if SortByStatus && OutputFormat == "text" {
		fmt.Println("\n===============[정렬된 결과]===============")
		printSortedTextResults()
		return
	}

	// HTML/PDF 출력에서 정렬 모드인 경우
	if SortByStatus && (OutputFormat == "html" || OutputFormat == "pdf") {
		processSortedHtmlResults()
	}

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

// printSortedTextResults 정렬된 결과를 텍스트로 출력
func printSortedTextResults() {
	// 결과를 상태별로 정렬
	sort.SliceStable(sortedResults, func(i, j int) bool {
		// 먼저 Pass가 먼저 오도록
		if sortedResults[i].Passed && !sortedResults[j].Passed {
			return true
		}
		// 다음으로 Fail이 오도록
		if !sortedResults[i].Passed && !sortedResults[i].Manual &&
			(sortedResults[j].Passed || sortedResults[j].Manual) {
			return true
		}
		// 마지막으로 Manual이 오도록
		if sortedResults[i].Manual && sortedResults[j].Passed {
			return false
		}
		if sortedResults[i].Manual && !sortedResults[j].Manual && !sortedResults[j].Passed {
			return false
		}

		// 같은 상태 내에서는 카테고리로 정렬
		if sortedResults[i].Category != sortedResults[j].Category {
			return sortedResults[i].Category < sortedResults[j].Category
		}

		// 마지막으로 체크 이름으로 정렬
		return sortedResults[i].CheckName < sortedResults[j].CheckName
	})

	// PASS 섹션 출력
	if hasPassed := countResults(sortedResults, true, false); hasPassed > 0 {
		fmt.Printf("\n===============[PASS]===============\n")
		for _, r := range sortedResults {
			if r.Passed {
				printSingleResult(r)
			}
		}
	}

	// FAIL 섹션 출력
	if hasFailed := countResults(sortedResults, false, false); hasFailed > 0 {
		fmt.Printf("\n===============[FAIL]===============\n")
		for _, r := range sortedResults {
			if !r.Passed && !r.Manual {
				printSingleResult(r)
			}
		}
	}

	// MANUAL 섹션 출력
	if hasManual := countResults(sortedResults, false, true); hasManual > 0 {
		fmt.Printf("\n===============[MANUAL]===============\n")
		for _, r := range sortedResults {
			if r.Manual {
				printSingleResult(r)
			}
		}
	}

	fmt.Println("\n===============[Checklist Summary]===============")
	fmt.Printf(Green+"✔ PASS: %d\n"+Reset, PassedCount)
	fmt.Printf(Red+"✖ FAIL: %d\n"+Reset, FailedCount)
	fmt.Printf(Yellow+"⚠ Manual: %d\n"+Reset, ManualCount)
	fmt.Println("===============[End of Summary]=================")
}

// countResults 특정 상태의 결과 개수 계산
func countResults(results []CheckResult, passed bool, manual bool) int {
	count := 0
	for _, r := range results {
		if r.Passed == passed && r.Manual == manual {
			count++
		}
	}
	return count
}

// processSortedHtmlResults HTML 출력을 위해 정렬된 결과 처리
func processSortedHtmlResults() {
	// HTML 출력용으로 모든 결과를 상태별로 변환
	for _, r := range sortedResults {
		status := "PASS"
		statusClass := "success"

		if !r.Passed {
			if r.Manual {
				status = "MANUAL"
				statusClass = "warning"
			} else {
				status = "FAIL"
				statusClass = "danger"
			}
		}

		htmlResult := CheckResultHTML{
			CheckName:   r.CheckName,
			Status:      status,
			StatusClass: statusClass,
			FailureMsg:  r.FailureMsg,
			Resources:   r.Resources,
			Runbook:     r.Runbook,
			Category:    r.Category,
		}

		sortedHtmlResults = append(sortedHtmlResults, htmlResult)

		// 원래 카테고리별로도 결과 추가 (템플릿이 카테고리 뷰를 지원하도록)
		if _, exists := categoryResults[r.Category]; !exists {
			categoryResults[r.Category] = []CheckResultHTML{}
			categoryOrder = append(categoryOrder, r.Category)
		}
		categoryResults[r.Category] = append(categoryResults[r.Category], htmlResult)
	}

	// 정렬된 결과를 htmlResults에 설정
	htmlResults = sortedHtmlResults

	// 상태별 카테고리 추가
	categoryResults["PASS"] = []CheckResultHTML{}
	categoryResults["FAIL"] = []CheckResultHTML{}
	categoryResults["MANUAL"] = []CheckResultHTML{}

	// categoryOrder 맨 앞에 상태 카테고리 추가
	categoryOrder = append([]string{"PASS", "FAIL", "MANUAL"}, categoryOrder...)

	// 각 상태별 결과 분류
	for _, r := range sortedHtmlResults {
		switch r.Status {
		case "PASS":
			categoryResults["PASS"] = append(categoryResults["PASS"], r)
		case "FAIL":
			categoryResults["FAIL"] = append(categoryResults["FAIL"], r)
		case "MANUAL":
			categoryResults["MANUAL"] = append(categoryResults["MANUAL"], r)
		}
	}
}

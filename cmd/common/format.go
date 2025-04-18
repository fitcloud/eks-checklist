package common

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// OutputFormat은 현재 출력 형식 설정을 저장합니다
var (
	OutputFormat string // "text", "html", "pdf"
)

// SetOutputFormat은 출력 형식을 설정
func SetOutputFormat(format string) {
	OutputFormat = strings.ToLower(format)
}

// HTMLTemplateData HTML 템플릿에 사용될 데이터 구조
type HTMLTemplateData struct {
	Title         string
	Date          string
	Results       []CheckResultHTML
	Summary       SummaryData
	Categories    map[string][]CheckResultHTML
	HasCategory   bool
	CategoryOrder []string
	SortByStatus  bool
}

// SummaryData 요약 데이터 구조
type SummaryData struct {
	PassCount   int
	FailCount   int
	ManualCount int
	Total       int
}

// 결과를 저장할 배열
var htmlResults []CheckResultHTML
var categoryResults map[string][]CheckResultHTML

// 카테고리 순서를 저장할 슬라이스 추가
var categoryOrder []string

// InitHTMLOutput HTML 출력 초기화
func InitHTMLOutput() {
	htmlResults = []CheckResultHTML{}
	categoryResults = make(map[string][]CheckResultHTML)
	categoryOrder = []string{} // 카테고리 순서 초기화
}

// AddResultForHTML HTML 출력을 위한 결과 추가
func AddResultForHTML(r CheckResult, category string) {
	// 맵이 nil인 경우 초기화
	if htmlResults == nil {
		htmlResults = []CheckResultHTML{}
	}
	if categoryResults == nil {
		categoryResults = make(map[string][]CheckResultHTML)
	}

	// 필터 기준에 따라 결과를 저장할지 확인
	if !ShouldPrintResult(r.Passed, r.Manual) {
		return
	}

	status := "PASS"
	statusClass := "success" // bootstrap 클래스에 맞게 수정

	if !r.Passed {
		if r.Manual {
			status = "MANUAL"
			statusClass = "warning" // bootstrap 경고 클래스
		} else {
			status = "FAIL"
			statusClass = "danger" // bootstrap 위험 클래스
		}
	}

	htmlResult := CheckResultHTML{
		CheckName:   r.CheckName,
		Status:      status,
		StatusClass: statusClass,
		FailureMsg:  r.FailureMsg,
		Resources:   r.Resources,
		Runbook:     r.Runbook,
		Category:    category,
	}

	htmlResults = append(htmlResults, htmlResult)

	// 카테고리별 결과 저장
	if _, exists := categoryResults[category]; !exists {
		categoryResults[category] = []CheckResultHTML{}
		categoryOrder = append(categoryOrder, category) // 카테고리 순서 기록
	}
	categoryResults[category] = append(categoryResults[category], htmlResult)
}

// SaveHTMLReport HTML 보고서 저장
func SaveHTMLReport() (string, error) {
	// 파일 생성
	now := time.Now()
	filename := "eks-checklist-report-" + now.Format("20060102-150405") + ".html"
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("파일 생성 오류: %v", err)
	}
	defer file.Close()

	// 템플릿 로드 - 외부 파일에서 템플릿 파싱
	tmpl, err := loadTemplate()
	if err != nil {
		return "", fmt.Errorf("템플릿 로딩 오류: %v", err)
	}

	// 템플릿 데이터 설정
	data := HTMLTemplateData{
		Title:   "EKS 체크리스트 결과 보고서",
		Date:    now.Format("2006-01-02 15:04:05"),
		Results: htmlResults,
		Summary: SummaryData{
			PassCount:   PassedCount,
			FailCount:   FailedCount,
			ManualCount: ManualCount,
			Total:       PassedCount + FailedCount + ManualCount,
		},
		Categories:    categoryResults,
		HasCategory:   len(categoryResults) > 0,
		CategoryOrder: categoryOrder,
		SortByStatus:  SortByStatus,
	}

	// 파일에 템플릿 실행 결과 저장
	err = tmpl.Execute(file, data)
	if err != nil {
		return "", fmt.Errorf("템플릿 실행 오류: %v", err)
	}

	return filename, nil
}

// loadTemplate 템플릿 파일 로드 함수
func loadTemplate() (*template.Template, error) {
	// 템플릿 파일 경로
	templatePath := "templates/report.html"

	// 템플릿 파일 존재 여부 확인
	if _, err := os.Stat(templatePath); err == nil {
		// 템플릿 파일이 존재하면 로드
		return template.ParseFiles(templatePath)
	} else if os.IsNotExist(err) {
		// 템플릿 파일이 없는 경우 프로젝트 루트 기준으로 다시 시도
		rootTemplatePath := filepath.Join("/workspaces/eks-checklist", templatePath)
		if _, err := os.Stat(rootTemplatePath); err == nil {
			return template.ParseFiles(rootTemplatePath)
		}

		// 현재 실행 경로 기준으로 시도
		execPath, err := os.Executable()
		if err == nil {
			execDir := filepath.Dir(execPath)
			execTemplatePath := filepath.Join(execDir, templatePath)
			if _, err := os.Stat(execTemplatePath); err == nil {
				return template.ParseFiles(execTemplatePath)
			}
		}

		// 다양한 경로를 시도해도 템플릿 파일을 찾지 못한 경우 에러 반환
		return nil, fmt.Errorf("HTML 템플릿 파일을 찾을 수 없습니다: %s", templatePath)
	} else {
		// 기타 오류
		return nil, fmt.Errorf("템플릿 파일 확인 중 오류 발생: %v", err)
	}
}

// ConvertHTMLToPDF HTML 보고서를 PDF로 변환
func ConvertHTMLToPDF(htmlFilePath string) (string, error) {
	// PDF 파일 이름 생성 (HTML 파일명에서 .html을 .pdf로 변경)
	pdfFilePath := strings.TrimSuffix(htmlFilePath, ".html") + ".pdf"

	f, err := os.Open(htmlFilePath)
	if err != nil {
		return "", err
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return "", err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(f))

	if err := pdfg.Create(); err != nil {
		return "", err
	}

	if err := pdfg.WriteFile(pdfFilePath); err != nil {
		return "", err
	}

	// 임시 HTML 파일 삭제
	if err := os.Remove(htmlFilePath); err != nil {
		fmt.Printf("경고: HTML 파일 삭제 실패: %v\n", err)
		// HTML 삭제 실패는 전체 프로세스의 실패로 간주하지 않음
	}

	return pdfFilePath, nil
}

package common

// 출력 형식 상수
const (
	OutputFormatText = "text"
	OutputFormatHTML = "html"
	OutputFormatPDF  = "pdf"
)

// 출력 필터 상수
const (
	OutputFilterAll    = "all"
	OutputFilterPass   = "pass"
	OutputFilterFail   = "fail"
	OutputFilterManual = "manual"
)

// 상태 표시 상수
const (
	StatusPass   = "PASS"
	StatusFail   = "FAIL"
	StatusManual = "MANUAL"
)

// CSS 클래스 상수
const (
	ClassSuccess = "success"
	ClassDanger  = "danger"
	ClassWarning = "warning"
)

// 색상 코드 상수
const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorReset  = "\033[0m"
)

// 카테고리 이름 상수
const (
	CategoryGeneral     = "General Check"
	CategorySecurity    = "Security Check"
	CategoryScalability = "Scalability Check"
	CategoryStability   = "Stability Check"
	CategoryNetwork     = "Network Check"
	CategoryCost        = "Cost Check"
)

// 파일 경로 상수
const (
	OutputDirPath     = "output/"
	TemplatesFilePath = "templates/report.html"
)

// 아이콘 상수
const (
	IconPass   = "✔"
	IconFail   = "✖"
	IconManual = "⚠"
)

// 유효한 출력 필터 목록 반환
func GetValidOutputFilters() []string {
	return []string{OutputFilterAll, OutputFilterPass, OutputFilterFail, OutputFilterManual}
}

// 유효한 출력 형식 목록 반환
func GetValidOutputFormats() []string {
	return []string{OutputFormatText, OutputFormatHTML, OutputFormatPDF}
}

// 모든 카테고리 목록 반환
func GetAllCategories() []string {
	return []string{
		CategoryGeneral,
		CategorySecurity,
		CategoryScalability,
		CategoryStability,
		CategoryNetwork,
		CategoryCost,
	}
} 
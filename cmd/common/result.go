package common

// CheckResult 체크 결과를 저장하는 구조체
type CheckResult struct {
	CheckName  string
	Passed     bool
	Manual     bool
	FailureMsg string
	Resources  []string
	Runbook    string
	Category   string // 카테고리 정보 추가
}

// CheckResultHTML HTML 출력을 위한 체크 결과 구조체
type CheckResultHTML struct {
	CheckName   string
	Status      string
	StatusClass string
	FailureMsg  string
	Resources   []string
	Runbook     string
	Category    string
}

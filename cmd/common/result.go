package common

type CheckResult struct {
	CheckName  string   // 체크 기능 이름?
	Manual     bool     // 수동 체크 여부
	Passed     bool     // 체크 통과 여부
	SuccessMsg string   // 성공 메세지
	FailureMsg string   // 실패 메세지
	Resources  []string // 영향받는 리소스 목록
	Runbook    string   // 런북 URL
}

package common

import (
	"strings"
)

// OutputFilter는 현재 출력 필터 설정을 저장합니다
var OutputFilter string

// SetOutputFilter는 출력 필터 타입을 설정합니다
func SetOutputFilter(filter string) {
	OutputFilter = strings.ToLower(filter)
}

// ShouldPrintResult는 필터에 따라 결과를 출력할지 결정합니다
func ShouldPrintResult(passed bool, manual bool) bool {
	if OutputFilter == "" {
		return true // 필터가 없으면 모든 결과 출력
	}

	resultType := ""
	if passed {
		resultType = "pass"
	} else if manual {
		resultType = "manual"
	} else {
		resultType = "fail"
	}

	return OutputFilter == resultType
}

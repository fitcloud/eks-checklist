package common

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintResult(t *testing.T) {
	// 테스트를 위해 표준 출력을 임시로 변경
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 테스트 후 원래 표준 출력으로 복원
	defer func() {
		os.Stdout = oldStdout
	}()

	// 테스트 케이스 설정
	testCases := []struct {
		name     string
		result   CheckResult
		expected []string
	}{
		{
			name: "Pass Result",
			result: CheckResult{
				CheckName: "Test Pass Check",
				Passed:    true,
				Manual:    false,
			},
			expected: []string{IconPass, "PASS", "Test Pass Check"},
		},
		{
			name: "Fail Result",
			result: CheckResult{
				CheckName:  "Test Fail Check",
				Passed:    false,
				Manual:    false,
				FailureMsg: "This is a failure message",
				Runbook:    "https://example.com/runbook",
			},
			expected: []string{IconFail, "FAIL", "Test Fail Check", "This is a failure message", "Runbook"},
		},
		{
			name: "Manual Result",
			result: CheckResult{
				CheckName:  "Test Manual Check",
				Passed:    false,
				Manual:    true,
				FailureMsg: "This needs manual check",
				Runbook:    "https://example.com/manual",
			},
			expected: []string{IconManual, "MANUAL", "Test Manual Check", "This needs manual check", "Runbook"},
		},
	}

	// 각 테스트 케이스 실행
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 테스트 전 카운터 초기화
			PassedCount = 0
			FailedCount = 0
			ManualCount = 0
			
			// 텍스트 모드로 설정
			OutputFormat = OutputFormatText
			SortByStatus = false
			
			// 결과 출력
			PrintResult(tc.result)
			
			// 파이프에서 출력 읽기
			w.Close()
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()
			
			// 예상 결과가 출력에 포함되어 있는지 확인
			for _, expected := range tc.expected {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain '%s', but got: %s", expected, output)
				}
			}
			
			// 카운터가 올바르게 증가했는지 확인
			if tc.result.Passed && PassedCount != 1 {
				t.Errorf("Expected PassedCount to be 1, but got %d", PassedCount)
			} else if !tc.result.Passed && tc.result.Manual && ManualCount != 1 {
				t.Errorf("Expected ManualCount to be 1, but got %d", ManualCount)
			} else if !tc.result.Passed && !tc.result.Manual && FailedCount != 1 {
				t.Errorf("Expected FailedCount to be 1, but got %d", FailedCount)
			}
			
			// 테스트를 위한 새 파이프 생성
			r, w, _ = os.Pipe()
			os.Stdout = w
		})
	}
}

func TestPrintCategoryHeader(t *testing.T) {
	// 테스트를 위해 표준 출력을 임시로 변경
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 테스트 후 원래 표준 출력으로 복원
	defer func() {
		os.Stdout = oldStdout
	}()

	// 테스트 케이스
	testCases := []struct {
		name       string
		category   string
		sortMode   bool
		outputFormat string
		shouldPrint bool
	}{
		{
			name:       "Print in text mode",
			category:   "Test Category",
			sortMode:   false,
			outputFormat: OutputFormatText,
			shouldPrint: true,
		},
		{
			name:       "Don't print in sort mode",
			category:   "Test Category",
			sortMode:   true,
			outputFormat: OutputFormatText,
			shouldPrint: false,
		},
		{
			name:       "Don't print in HTML mode",
			category:   "Test Category",
			sortMode:   false,
			outputFormat: OutputFormatHTML,
			shouldPrint: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 설정
			SortByStatus = tc.sortMode
			OutputFormat = tc.outputFormat
			
			// 헤더 출력
			PrintCategoryHeader(tc.category)
			
			// 파이프에서 출력 읽기
			w.Close()
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()
			
			// 예상 결과 확인
			expectedHeader := "===============[" + tc.category + "]==============="
			if tc.shouldPrint && !strings.Contains(output, expectedHeader) {
				t.Errorf("Expected output to contain '%s', but got: %s", expectedHeader, output)
			} else if !tc.shouldPrint && strings.Contains(output, expectedHeader) {
				t.Errorf("Expected output not to contain '%s', but got: %s", expectedHeader, output)
			}
			
			// 현재 카테고리가 올바르게 설정되었는지 확인
			if CurrentCategory != tc.category {
				t.Errorf("Expected CurrentCategory to be '%s', but got '%s'", tc.category, CurrentCategory)
			}
			
			// 테스트를 위한 새 파이프 생성
			r, w, _ = os.Pipe()
			os.Stdout = w
		})
	}
} 
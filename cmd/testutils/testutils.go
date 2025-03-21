package testutils

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

// LoadTestCases는 주어진 파일 경로에서 YAML 형식의 테스트 데이터를 로드합니다.
func LoadTestCases(t *testing.T, filePath string) []map[string]interface{} {
	// testdata 디렉토리 내에서 상대 경로로 파일을 찾도록 설정
	absPath, err := filepath.Abs(filepath.Join("..", "..", "testdata", filePath))
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	// 파일을 읽어옵니다.
	data, err := os.ReadFile(absPath)
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}

	// YAML 파싱
	var cases []map[string]interface{}
	if err := yaml.Unmarshal(data, &cases); err != nil {
		t.Fatalf("failed to parse YAML: %v", err)
	}

	return cases
}

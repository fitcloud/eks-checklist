package common

import (
	"encoding/json"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

// Kubernetes 리소스를 YAML로 저장
func SaveK8sResourceAsYAML(obj runtime.Object, path string) error {
	serializer := k8sjson.NewSerializerWithOptions(
		k8sjson.DefaultMetaFactory, nil, nil,
		k8sjson.SerializerOptions{Yaml: true, Pretty: true, Strict: false},
	)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("파일 생성 실패: %w", err)
	}
	defer f.Close()

	if err := serializer.Encode(obj, f); err != nil {
		return fmt.Errorf("YAML 인코딩 실패: %w", err)
	}

	return nil
}

// 일반 데이터를 JSON으로 저장
func SaveAsJSON(data interface{}, path string) error {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON 마샬 실패: %w", err)
	}

	if err := os.WriteFile(path, jsonBytes, 0644); err != nil {
		return fmt.Errorf("파일 저장 실패: %w", err)
	}

	return nil
}

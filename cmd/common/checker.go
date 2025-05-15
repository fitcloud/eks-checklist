package common

import (
	"fmt"
	"sort"
)

// Checker 인터페이스는 체크 항목을 수행하는 모든 함수가 구현해야 하는 인터페이스입니다.
type Checker interface {
	// Check 함수는 체크를 수행하고 결과를 반환합니다.
	Check() CheckResult
	
	// GetName 함수는 체크 항목의 이름을 반환합니다.
	GetName() string
	
	// GetCategory 함수는 체크 항목의 카테고리를 반환합니다.
	GetCategory() string
}

// CheckerFunc는 함수를 Checker 인터페이스로 변환하기 위한 어댑터 타입입니다.
type CheckerFunc struct {
	Name     string
	Category string
	Func     func() CheckResult
}

// Check 함수는 CheckerFunc가 Checker 인터페이스를 구현하도록 합니다.
func (cf CheckerFunc) Check() CheckResult {
	return cf.Func()
}

// GetName 함수는 체크 항목의 이름을 반환합니다.
func (cf CheckerFunc) GetName() string {
	return cf.Name
}

// GetCategory 함수는 체크 항목의 카테고리를 반환합니다.
func (cf CheckerFunc) GetCategory() string {
	return cf.Category
}

// CheckerRegistry는 등록된 모든 체크 항목을 관리합니다.
type CheckerRegistry struct {
	checkers map[string][]Checker // 카테고리별 체크 항목 맵
}

// NewCheckerRegistry는 새로운 CheckerRegistry를 생성합니다.
func NewCheckerRegistry() *CheckerRegistry {
	return &CheckerRegistry{
		checkers: make(map[string][]Checker),
	}
}

// Register는 새로운 체크 항목을 등록합니다.
func (r *CheckerRegistry) Register(checker Checker) {
	category := checker.GetCategory()
	if _, exists := r.checkers[category]; !exists {
		r.checkers[category] = []Checker{}
	}
	r.checkers[category] = append(r.checkers[category], checker)
}

// RegisterFunc는 함수를 체크 항목으로 등록합니다.
func (r *CheckerRegistry) RegisterFunc(name, category string, fn func() CheckResult) {
	r.Register(CheckerFunc{
		Name:     name,
		Category: category,
		Func:     fn,
	})
}

// GetCheckers는 특정 카테고리의 모든 체크 항목을 반환합니다.
func (r *CheckerRegistry) GetCheckers(category string) []Checker {
	if checkers, exists := r.checkers[category]; exists {
		return checkers
	}
	return []Checker{}
}

// GetCategories는 등록된 모든 카테고리를 반환합니다.
func (r *CheckerRegistry) GetCategories() []string {
	categories := make([]string, 0, len(r.checkers))
	for category := range r.checkers {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	return categories
}

// RunChecks는 모든 체크를 실행하고 결과를 출력합니다.
func (r *CheckerRegistry) RunChecks() {
	categories := r.GetCategories()
	for _, category := range categories {
		PrintCategoryHeader(category)
		checkers := r.GetCheckers(category)
		for _, checker := range checkers {
			result := checker.Check()
			PrintResult(result)
		}
	}
	PrintSummary()
}

// RunCategoryChecks는 특정 카테고리의 모든 체크를 실행하고 결과를 출력합니다.
func (r *CheckerRegistry) RunCategoryChecks(category string) {
	PrintCategoryHeader(category)
	checkers := r.GetCheckers(category)
	for _, checker := range checkers {
		result := checker.Check()
		PrintResult(result)
	}
}

// DefaultRegistry는 기본 체크 항목 레지스트리입니다.
var DefaultRegistry = NewCheckerRegistry() 
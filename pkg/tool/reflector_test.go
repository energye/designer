package tool

import (
	"testing"
)

type TestStruct struct{}

func (t TestStruct) Method1()  {}
func (t *TestStruct) Method2() {}
func (t TestStruct) Method3()  {}

type TestStruct2 struct {
	TestStruct
}

func (i TestStruct2) Method21()  {}
func (i *TestStruct2) Method22() {}

type TestStruct3 struct {
	TestStruct2
}

func (i TestStruct3) Method31()  {}
func (i *TestStruct3) Method32() {}

func TestSimpleObjectMethodNames(t *testing.T) {
	var ts3 TestStruct3
	methods := GetObjectMethodNames(&ts3)
	t.Log("方法列表:")
	methods.Iterate(func(key string, method *TMethod) bool {
		t.Log("方法名:", method.Name, "所属结构体:", method.StructName, "层级:", method.Level)
		return false
	})
}

// 测试GetObjectMethodNames函数
func TestGetObjectMethodNames(t *testing.T) {
	// 测试用例定义
	tests := []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name:  "结构体值类型测试",
			input: TestStruct{},
			expected: []string{
				"Method1",
				"Method2",
				"Method3",
			},
		},
		{
			name:  "结构体指针类型测试",
			input: &TestStruct{},
			expected: []string{
				"Method1",
				"Method2",
				"Method3",
			},
		},
		{
			name:  "结构体值类型测试",
			input: TestStruct2{},
			expected: []string{
				"Method1",
				"Method2",
				"Method3",
				"Method21",
				"Method22",
			},
		},
		{
			name:  "结构体指针类型测试",
			input: &TestStruct2{},
			expected: []string{
				"Method1",
				"Method2",
				"Method3",
				"Method21",
				"Method22",
			},
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetObjectMethodNames(tt.input)

			// 验证返回值不为nil
			if result == nil {
				t.Errorf("GetObjectMethodNames() = nil, want non-nil ArrayMap")
				return
			}

			// 验证方法数量
			if len(result.Keys()) != len(tt.expected) {
				t.Errorf("GetObjectMethodNames() returned %d methods, want %d", len(result.Keys()), len(tt.expected))
				return
			}

			// 验证每个期望的方法都存在
			for _, expectedMethod := range tt.expected {
				found := false
				if method := result.Get(expectedMethod); method.Name != "" {
					found = true
					break
				}
				if !found {
					t.Errorf("Expected method %s not found in result", expectedMethod)
				}
			}
		})
	}
}

// TestGetObjectMethodNames_EmptyStruct 测试空结构体
func TestGetObjectMethodNames_EmptyStruct(t *testing.T) {
	type EmptyStruct struct{}

	result := GetObjectMethodNames(EmptyStruct{})

	if result == nil {
		t.Errorf("GetObjectMethodNames() = nil, want non-nil ArrayMap")
		return
	}

	if len(result.Keys()) != 0 {
		t.Errorf("GetObjectMethodNames() returned %d methods for empty struct, want 0", len(result.Keys()))
	}
}

// TestGetObjectMethodNames_NilPointer 测试nil指针
func TestGetObjectMethodNames_NilPointer(t *testing.T) {
	var nilPtr *TestStruct

	result := GetObjectMethodNames(nilPtr)

	if result == nil {
		t.Errorf("GetObjectMethodNames() = nil, want non-nil ArrayMap")
		return
	}

	// nil指针应该能获取到值接收者的方法
	expectedMethods := []string{"Method1", "Method2"}
	if len(result.Keys()) != len(expectedMethods) {
		t.Errorf("GetObjectMethodNames() returned %d methods, want %d", len(result.Keys()), len(expectedMethods))
	}
}

// BenchmarkGetObjectMethodNames 性能测试
func BenchmarkGetObjectMethodNames(b *testing.B) {
	testStruct := TestStruct{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetObjectMethodNames(testStruct)
	}
}

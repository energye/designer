package codegen

import (
	"encoding/json"
	"fmt"
	"github.com/energye/designer/uigen"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// go 代码生成 自动时时生成
// 依赖 uigen UI 布局文件
// 生成触发条件: 检查文件修改

// GenerateCode 根据UI文件生成Go代码
func GenerateCode(uiFilePath string) error {
	// 读取并解析UI文件
	data, err := os.ReadFile(uiFilePath)
	if err != nil {
		return fmt.Errorf("读取UI文件失败: %w", err)
	}

	var uiComponent uigen.UIComponent
	if err := json.Unmarshal(data, &uiComponent); err != nil {
		return fmt.Errorf("解析UI文件失败: %w", err)
	}

	// 生成自动代码文件 (main_form.ui.go)
	if err := generateAutoCode(uiFilePath, uiComponent); err != nil {
		return fmt.Errorf("生成自动代码失败: %w", err)
	}

	// 生成用户代码文件 (main_form.go) - 仅当文件不存在时
	if err := generateUserCode(uiFilePath, uiComponent); err != nil {
		return fmt.Errorf("生成用户代码失败: %w", err)
	}

	return nil
}

// generateAutoCode 生成自动代码文件
func generateAutoCode(uiFilePath string, component uigen.UIComponent) error {
	// 构建模板数据
	data := buildAutoTemplateData(component)

	// 解析模板
	tmpl, err := template.New("auto").Parse(autoCodeTemplate)
	if err != nil {
		return fmt.Errorf("解析自动代码模板失败: %w", err)
	}

	// 生成代码
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("执行自动代码模板失败: %w", err)
	}

	// 格式化代码
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		return fmt.Errorf("格式化代码失败: %w", err)
	}

	// 写入文件
	baseName := strings.TrimSuffix(filepath.Base(uiFilePath), filepath.Ext(uiFilePath))
	autoFileName := baseName + ".ui.go"
	autoFilePath := filepath.Join(filepath.Dir(uiFilePath), autoFileName)

	if err := os.WriteFile(autoFilePath, formatted, 0644); err != nil {
		return fmt.Errorf("写入自动代码文件失败: %w", err)
	}

	return nil
}

// generateUserCode 生成用户代码文件
func generateUserCode(uiFilePath string, component uigen.UIComponent) error {
	// 检查文件是否已存在
	baseName := strings.TrimSuffix(filepath.Base(uiFilePath), filepath.Ext(uiFilePath))
	userFileName := baseName + ".go"
	userFilePath := filepath.Join(filepath.Dir(uiFilePath), userFileName)

	// 如果文件已存在，不覆盖
	if _, err := os.Stat(userFilePath); err == nil {
		return nil // 文件已存在，直接返回
	}

	// 构建模板数据
	data := buildUserTemplateData(component)

	// 解析模板
	tmpl, err := template.New("user").Parse(userCodeTemplate)
	if err != nil {
		return fmt.Errorf("解析用户代码模板失败: %w", err)
	}

	// 生成代码
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("执行用户代码模板失败: %w", err)
	}

	// 格式化代码
	formatted, err := format.Source([]byte(buf.String()))
	if err != nil {
		return fmt.Errorf("格式化代码失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(userFilePath, formatted, 0644); err != nil {
		return fmt.Errorf("写入用户代码文件失败: %w", err)
	}

	return nil
}

package mappergen

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// TypeInfo 存储类型定义信息
type TypeInfo struct {
	Name       string // 类型名（如TAlign）
	Underlying string // 底层类型（如int32）
}

// ConstInfo 存储常量信息
type ConstInfo struct {
	Name  string // 常量名（如AlTop）
	Type  string // 常量类型（如TAlign，可能为空）
	Value string // 常量值（如"0"、"iota"）
}

// LCL 的类型映射
func LCLMapper() {
	// 保存目录
	wd, _ := os.Getwd()
	outPath := filepath.Join(wd, "pkg", "mapper", "lcl_types.go")
	pkgPath := "C:\\app\\workspace\\gen\\gout\\lcl\\go\\types"

	// 扫描包并获取类型和常量信息
	types, consts, err := scanPackage(pkgPath)
	if err != nil {
		fmt.Printf("扫描失败: %v\n", err)
		os.Exit(1)
	}
	_ = types //
	// 打印结果
	//fmt.Println("===== 类型定义 =====")
	//for _, t := range types {
	//	fmt.Printf("类型名: %-20s 底层类型: %s\n", t.Name, t.Underlying)
	//}
	//
	//fmt.Println("\n===== 常量定义 =====")
	//for _, c := range consts {
	//	if c.Type != "" {
	//		fmt.Printf("常量名: %-20s 类型: %-15s 值: %s\n", c.Name, c.Type, c.Value)
	//	} else {
	//		fmt.Printf("常量名: %-20s 值: %s\n", c.Name, c.Value)
	//	}
	//}
	mapperTemplate := `package mapper

import . "github.com/energye/lcl/types"

var lclTypesMapper = make(map[string]any)

func init() {
{{mappers}}
}

// 获取映射的类型值
func GetLCL(name string) any {
	return lclTypesMapper[name]
}
`
	constBuf := bytes.Buffer{}
	for _, c := range consts {
		// typesMapper["AlClient"] = AlClient
		constBuf.WriteString("\tlclTypesMapper[")
		constBuf.WriteString(`"` + c.Name + `"`)
		constBuf.WriteString("] = ")
		constBuf.WriteString(c.Name)
		constBuf.WriteString("\n")
	}
	// 保存
	mapperTemplate = strings.Replace(mapperTemplate, "{{mappers}}", constBuf.String(), 1)
	os.WriteFile(outPath, []byte(mapperTemplate), fs.ModePerm)
}

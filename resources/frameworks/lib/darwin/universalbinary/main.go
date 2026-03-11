package main

import (
	"fmt"
	"github.com/energye/designer/pkg/tool"
	"github.com/energye/designer/resources/frameworks/lib"
	"github.com/energye/lcl/rtl/version"
	"os"
	"path/filepath"
)

func main() {
	version.Init()
	wd, _ := os.Getwd()
	output := filepath.Join(wd, "resources", "frameworks", "lib")
	fmt.Println("UniversalBinary output:", output)
	amd64ZipPath := filepath.Join(output, lib.PathAMD64Cocoa)
	arm64ZipPath := filepath.Join(output, lib.PathARM64Cocoa)
	outputPath := filepath.Join(output, "darwin")
	universalLibFilePath, err := tool.UniversalBinary(amd64ZipPath, arm64ZipPath, outputPath)
	fmt.Println("UniversalBinary universalLibFilePath:", universalLibFilePath, "error:", err)
}

package main

import (
	"github.com/energye/designer/designer"
	_ "github.com/energye/designer/pkg/syso"
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/lcl"
	"os"
	"path/filepath"
)

func main() {
	libname.LibName = func() string {
		wd, _ := os.Getwd()
		return filepath.Join(wd, "../", "gen", "gout", "liblcl.dll")
	}()
	lcl.Init(nil, nil)
	designer.Run()
}

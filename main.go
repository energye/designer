package main

import (
	"github.com/energye/designer/designer"
	_ "github.com/energye/designer/pkg/syso"
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/lcl"
)

func main() {
	libname.LibName = "E:\\SWT\\gopath\\src\\github.com\\energye\\workspace\\gen\\gout\\liblcl.dll"
	lcl.Init(nil, nil)
	designer.Run()
}

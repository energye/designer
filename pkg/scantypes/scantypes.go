package main

import "github.com/energye/designer/pkg/scantypes/mappergen"

// 生成类型映射.go
// 用于在动态设置属性时

func main() {
	mappergen.LCLMapper()
	mappergen.CEFMapper()
	mappergen.WV2Mapper()
}

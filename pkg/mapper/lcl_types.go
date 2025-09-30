package mapper

import . "github.com/energye/lcl/types"

var typesMapper = make(map[string]any)

func init() {
	typesMapper["AlClient"] = AlClient
}

// 获取映射的类型值
func Get(name string) any {
	return typesMapper[name]
}

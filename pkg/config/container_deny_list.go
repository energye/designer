package config

import (
	"encoding/json"
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/resources"
)

type containerDenyList map[string]struct{}

// 容器拒绝列表, 维护哪些组件不能做为容器
var ContainerDenyList = make(containerDenyList)

func init() {
	data := resources.ContainerDenyList()
	var dataList []string
	err.CheckErr(json.Unmarshal(data, &dataList))
	for _, componentName := range dataList {
		ContainerDenyList[componentName] = struct{}{}
	}
}

// 判断组件是否被配置为非容器
func (m containerDenyList) IsDeny(name string) bool {
	_, ok := m[name]
	return ok
}

package config

import (
	"bytes"
	"encoding/json"
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/designer/resources"
)

// 弹窗过滤
var DialogFilter = &dialogFilter{}

type dialogFilter struct {
	Image []string `json:"image"`
	File  []string `json:"file"`
}

func init() {
	data := resources.DialogFilter()
	if data == nil {
		logs.Error("加载弹窗过滤配置为nil")
		return
	}
	err := json.Unmarshal(data, DialogFilter)
	if data == nil {
		logs.Error("加载弹窗过滤配置错误:", err.Error())
		return
	}
}

func (m *dialogFilter) ImageFilter() string {
	buf := bytes.Buffer{}
	for i, item := range m.Image {
		if i > 0 {
			buf.WriteString("|")
		}
		buf.WriteString(item)
	}
	return buf.String()
}

func (m *dialogFilter) FileFilter() string {
	buf := bytes.Buffer{}
	for i, item := range m.File {
		if i > 0 {
			buf.WriteString("|")
		}
		buf.WriteString(item)
	}
	return buf.String()
}

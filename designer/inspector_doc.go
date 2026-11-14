// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package designer

import (
	"strings"

	"github.com/energye/designer/pkg/vtedit"
)

// 对象检查器 属性文档提示
// 悬停属性/事件节点时, 在检查器底部文档面板显示属性元数据信息

var gInspectorDocLastKey string // 上次显示的文档标识, 避免重复刷新

// InspectorShowPropertyDoc 悬停属性节点时更新文档面板
func InspectorShowPropertyDoc(data *vtedit.TEditNodeData) {
	if MainWindow.contentLayout == nil || MainWindow.contentLayout.layoutInspector == nil {
		return
	}
	m := MainWindow.contentLayout.layoutInspector
	if m.doc == nil || !m.doc.IsValid() {
		return
	}
	if data == nil || data.EditNodeData == nil {
		return
	}
	key := data.Name()
	if key == gInspectorDocLastKey {
		return
	}
	gInspectorDocLastKey = key
	m.doc.SetText(inspectorDocText(data))
}

func inspectorDocText(data *vtedit.TEditNodeData) string {
	var b strings.Builder
	b.WriteString("名称: ")
	b.WriteString(data.Name())
	md := data.EditNodeData.Metadata
	if md != nil {
		if md.Kind == "tkMethod" {
			b.WriteString("\n事件: 绑定 Go 事件处理函数名")
		} else if md.Kind != "" {
			b.WriteString("\n类型: ")
			b.WriteString(md.Kind)
			if md.Type != "" && md.Type != md.Kind {
				b.WriteString("  (")
				b.WriteString(md.Type)
				b.WriteString(")")
			}
		}
		if md.Options != "" {
			b.WriteString("\n可选值: ")
			b.WriteString(md.Options)
		}
	}
	return b.String()
}

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

package preview

import (
	"github.com/energye/designer/consts"
	"github.com/energye/designer/event"
	"github.com/energye/designer/pkg/logs"
)

func init() {
	event.On(event.Preview, func(trigger event.TTrigger) {
		ps, ok := trigger.Payload.(consts.PreviewState)
		if ok {
			if ps == consts.PsStarted {
				// 启动运行预览
				go runPreview(trigger.Result)
			} else if ps == consts.PsStop {
				// 停止运行预览
				go stopPreview()
			}
		} else {
			logs.Error("运行预览错误, 操作参数不正确, option:", ps)
		}
	}, func() {
		logs.Info("停止预览处理器")
	})
}

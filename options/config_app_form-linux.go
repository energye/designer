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

package options

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/lcl/lcl"
)

func (m *TConfigProjectForm) initLinuxOptions() {
	logs.Debug("TConfigProjectForm initLinuxOptions")
	test := lcl.NewLabel(m)
	test.SetLeft(m.platformTabPageLinux.Width()/2 - 20)
	test.SetTop(m.platformTabPageLinux.Height()/2 - 20)
	test.SetParent(m.platformTabPageLinux)
	test.SetCaption("Empty")
}

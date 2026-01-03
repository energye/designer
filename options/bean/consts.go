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

package bean

type MacOSUIElementList int32

const (
	MacOSUIElementListNo  MacOSUIElementList = iota // false (常规前台应用)
	MacOSUIElementListYes                           // true (后台应用, 无 Dock 图标)
)

type LSMinimumSystemVersion int32

const (
	LSMinimumSystemVersion_10_15 LSMinimumSystemVersion = iota // 10.15 (Intel)
	LSMinimumSystemVersion_11_0                                // 11.0 (Apple Silicon)
)

type GUIRenderFramework = string

const (
	GUIRenderFramework_LCL GUIRenderFramework = "LCL"
	GUIRenderFramework_WV  GUIRenderFramework = "WV"
	GUIRenderFramework_CEF GUIRenderFramework = "CEF"
)

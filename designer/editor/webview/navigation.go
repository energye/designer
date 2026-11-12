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

package webview

import (
	"fmt"
	"net/url"
	"strings"
)

type urlAction int

const (
	// urlActionInternal 留在当前 WebView 内导航
	urlActionInternal urlAction = iota
	// urlActionExternal 取消导航并交给系统默认程序打开
	urlActionExternal
	// urlActionBlock 取消导航且不向外打开
	urlActionBlock
)

func monacoAssetsOrigin(port int) string {
	if port <= 0 {
		return "http://127.0.0.1:22022"
	}
	return fmt.Sprintf("http://127.0.0.1:%d", port)
}

// classifyURL 仅同源 http(s)、about: 页面视为应用内导航，其余 http(s)/mailto 交给系统。
func classifyURL(rawURL, origin string) urlAction {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return urlActionBlock
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return urlActionBlock
	}
	switch strings.ToLower(u.Scheme) {
	case "about", "chrome-error":
		return urlActionInternal
	case "http", "https":
		if u.Host == "" {
			return urlActionBlock
		}
		o, err := url.Parse(origin)
		if err == nil && strings.EqualFold(u.Scheme, o.Scheme) && strings.EqualFold(u.Host, o.Host) {
			return urlActionInternal
		}
		return urlActionExternal
	case "mailto":
		return urlActionExternal
	default:
		return urlActionBlock
	}
}

// decideNavigation 返回是否取消顶级导航、是否交给系统默认程序打开。
func decideNavigation(rawURL, origin string) (cancel, openExternal bool) {
	switch classifyURL(rawURL, origin) {
	case urlActionInternal:
		return false, false
	case urlActionExternal:
		return true, true
	default:
		return true, false
	}
}

// decideNewWindow 新窗口一律不创建，仅外部 URL 交给系统默认程序。
func decideNewWindow(rawURL, origin string) (openExternal bool) {
	return classifyURL(rawURL, origin) == urlActionExternal
}

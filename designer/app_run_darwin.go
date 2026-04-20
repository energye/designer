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

//go:build darwin

package designer

import (
	"github.com/energye/designer/pkg/logs"
	"github.com/energye/energy/v3/platform/darwin/cocoa"
	"net/url"
	"strings"
)

func beforeRun() {
	cocoa.NSApp.InitAppDelegate()
	cocoa.NSApp.SetOnOpenURLs(func(urls []string) {
		logs.Info("SetOnOpenURLs:", strings.Join(urls, " "))
		if urls != nil && len(urls) > 0 {
			openUrl := urls[0]
			pUrl, err := url.Parse(openUrl)
			if err != nil {
				logs.Error("SetOnOpenURLs", err.Error())
				return
			}
			if pUrl.Scheme == "file" {
				// 文件关联 file:///Users/yanghy/app/lazdemo/myapp/myapp.egp
				associateFile = pUrl.Path
				stopAutoAssociateProjectLoad()
				loadProject(OpenFile())
			} else if pUrl.Scheme == "http" || pUrl.Scheme == "https" {
				// 通用链接 https://example.com/action?id=1
				fullLink := pUrl.Path
				if pUrl.RawQuery != "" {
					fullLink += "?" + pUrl.RawQuery
				}
				universalLink = fullLink
			} else {
				// 自定义协议 energy://open?id=123
				fullLink := pUrl.Path
				if pUrl.RawQuery != "" {
					fullLink += "?" + pUrl.RawQuery
				}
				associateProtocol = fullLink
			}
		}
	})
	cocoa.NSApp.SetOnUniversalLink(func(link string) {
		logs.Info("SetOnUniversalLink:", link)
		pUrl, err := url.Parse(link)
		if err != nil {
			return
		}
		if pUrl.Scheme == "http" || pUrl.Scheme == "https" {
			// 通用链接 https://example.com/action?id=1
			fullLink := pUrl.Path
			if pUrl.RawQuery != "" {
				fullLink += "?" + pUrl.RawQuery
			}
			universalLink = fullLink
		}
	})
}

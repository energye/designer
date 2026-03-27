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

//go:build windows || liball

package lib

import "embed"

var (
	//go:embed windows/MicrosoftEdgeWebview2Setup.zip
	webview2Setup embed.FS
)

func init() {
	Add(&EmbedFS{Path: PathWV2Setup, Lib: &webview2Setup, OutputFilename: "MicrosoftEdgeWebview2Setup.exe", NotReleased: true})
}

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
	//go:embed windows/libenergy-amd64.zip
	libAMD64Win32 embed.FS
	//go:embed windows/WebView2Loader-amd64.zip
	libWV2AMD64Win32 embed.FS
)

func init() {
	Add(&EmbedFS{Path: PathAMD64Win32, Lib: &libAMD64Win32, OutputFilename: "libenergy-amd64.dll"})
	Add(&EmbedFS{Path: PathWV2AMD64Win32, Lib: &libWV2AMD64Win32, OutputFilename: "WebView2Loader-amd64.dll"})
}

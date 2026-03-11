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

//go:build (windows && 386) || liball

package lib

var (
	//go:embed windows/libenergy-windows-i386.zip
	libI386Win32 embed.FS
	//go:embed windows/WebView2Loader-i386.zip
	libWV2I386Win32 embed.FS
)

func init() {
	libs.Add(PathI386Win32, &EmbedFS{Lib: &libI386Win32, OutputFilename: "libenergy-windows-i386-win32.dll"})
	libs.Add(PathWV2I386Win32, &EmbedFS{Lib: &libWV2I386Win32, OutputFilename: "WebView2Loader-i386.dll"})
}

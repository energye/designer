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

//go:build (linux && arm64) || liball

package lib

import "embed"

var (
	//go:embed linux/libenergy-linux-arm64-gtk2.zip
	libARM64GTK2 embed.FS
	//go:embed linux/libenergy-linux-arm64-gtk3.zip
	libARM64GTK3 embed.FS
)

const (
	pathARM64Gtk2 = "linux/libenergy-linux-arm64-gtk2.zip"
	pathARM64Gtk3 = "linux/libenergy-linux-arm64-gtk3.zip"
)

func init() {
	libs.Add(pathARM64Gtk2, &EmbedFS{Lib: &libARM64GTK2, OutputFilename: "libenergy-linux-arm64-gtk2.so"})
	libs.Add(pathARM64Gtk3, &EmbedFS{Lib: &libARM64GTK3, OutputFilename: "libenergy-linux-arm64-gtk3.so"})
}

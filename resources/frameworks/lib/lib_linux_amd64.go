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

//go:build (linux && amd64) || liball

package lib

import (
	"embed"
)

var (
	//go:embed linux/libenergy-linux-amd64-gtk2.zip
	libAMD64GTK2 embed.FS
	//go:embed linux/libenergy-linux-amd64-gtk3.zip
	libAMD64GTK3 embed.FS
)

const (
	pathAMD64Gtk2 = "linux/libenergy-linux-amd64-gtk2.zip"
	pathAMD64Gtk3 = "linux/libenergy-linux-amd64-gtk3.zip"
)

func init() {
	libs.Add(pathAMD64Gtk2, &EmbedFS{Lib: &libAMD64GTK2, OutputFilename: "libenergy-amd64-gtk2.so"})
	libs.Add(pathAMD64Gtk3, &EmbedFS{Lib: &libAMD64GTK3, OutputFilename: "libenergy-amd64-gtk3.so"})
}

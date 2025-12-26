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

//go:build (linux && arm) || liball

package lib

import "embed"

var (
	//go:embed linux/libenergy-linux-armhf-gtk2.zip
	libARMGTK2 embed.FS
	//go:embed linux/libenergy-linux-armhf-gtk3.zip
	libARMGTK3 embed.FS
)

const (
	pathARMGtk2 = "linux/libenergy-linux-armhf-gtk2.zip"
	pathARMGtk3 = "linux/libenergy-linux-armhf-gtk3.zip"
)

func init() {
	libs.Add(pathARMGtk2, &EmbedFS{Lib: &libARMGTK2, OutputFilename: "libenergy-linux-arm-gtk2.so"})
	libs.Add(pathARMGtk3, &EmbedFS{Lib: &libARMGTK3, OutputFilename: "libenergy-linux-arm-gtk3.so"})
}

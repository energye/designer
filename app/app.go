// Copyright © yanghy. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.

package app

import (
	"github.com/energye/energy/v3/application/pack"
	"github.com/energye/lcl/api/libname"
	"github.com/energye/lcl/lcl"
	"runtime"
)

var (
	// Info app pack info
	Info = pack.Info
)

// Forms form maintenance list
var Forms = []lcl.IEngForm{}

func init() {
	if runtime.GOOS == "linux" {
		libname.UseWS = "gtk3"
	}
}

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

package lib

import (
	"github.com/energye/designer/pkg/err"
	"github.com/energye/designer/pkg/logs"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	Path string
)

func init() {
	tempDir := os.TempDir()
	outPath := filepath.Join(tempDir, Name)
	libByte, e := lib.ReadFile(path)
	err.CheckErr(e)
	os.WriteFile(outPath, libByte, fs.ModePerm)
	Path = outPath
	logs.Info("Lib Path:", outPath)
}

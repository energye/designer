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

package project

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/energye/designer/pkg/config"
)

type GoModUpdateOptions struct {
	Dir      string
	OnOutput func(string)
}

func UpdateGoModDependencies(ctx context.Context, options GoModUpdateOptions) error {
	if strings.TrimSpace(options.Dir) == "" {
		return errors.New("project directory is empty")
	}
	writeGoModOutput(options.OnOutput, "Updating go mod. Please wait...")
	for modName, modVersion := range config.DesignerConfig.Dependencies {
		mod := modName + "@" + modVersion
		if err := runGoModCommand(ctx, options.Dir, options.OnOutput, "get", mod); err != nil {
			return err
		}
	}
	if err := runGoModCommand(ctx, options.Dir, options.OnOutput, "mod", "tidy"); err != nil {
		return err
	}
	writeGoModOutput(options.OnOutput, "Go mod update successfully")
	return nil
}

func runGoModCommand(ctx context.Context, dir string, onOutput func(string), args ...string) error {
	writeGoModOutput(onOutput, "go "+strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		writeGoModOutput(onOutput, strings.TrimSpace(string(output)))
	}
	if err != nil {
		return fmt.Errorf("go %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

func writeGoModOutput(onOutput func(string), message string) {
	if onOutput != nil && message != "" {
		onOutput(message)
	}
}

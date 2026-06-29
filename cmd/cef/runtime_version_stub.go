//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

//go:build !darwin && !linux && !windows

package cef

import "fmt"

func CEFLibVersion(path string) (RuntimeVersion, error) {
	return RuntimeVersion{}, fmt.Errorf("CEFLibVersion unsupported on this platform: %s", path)
}

//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

//go:build darwin || linux

package cef

import (
	"fmt"

	"github.com/ebitengine/purego"
)

func CEFLibVersion(path string) (RuntimeVersion, error) {
	handle, err := purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_LOCAL)
	if err != nil {
		return RuntimeVersion{}, err
	}
	defer purego.Dlclose(handle)
	var fn func(*int32, *int32, *int32, *int32)
	func() {
		defer func() {
			_ = recover()
		}()
		purego.RegisterLibFunc(&fn, handle, "CEFLibVersion")
	}()
	if fn == nil {
		return RuntimeVersion{}, fmt.Errorf("CEFLibVersion symbol not found in %s", path)
	}
	var version RuntimeVersion
	fn(&version.Major, &version.Minor, &version.Release, &version.Build)
	return version, nil
}

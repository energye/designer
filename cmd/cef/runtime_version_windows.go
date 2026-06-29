//----------------------------------------
//
// Copyright © yanghy. All Rights Reserved.
//
// Licensed under Apache License Version 2.0, January 2004
//
// https://www.apache.org/licenses/LICENSE-2.0
//
//----------------------------------------

//go:build windows

package cef

import (
	"fmt"
	"syscall"

	"github.com/ebitengine/purego"
)

func CEFLibVersion(path string) (RuntimeVersion, error) {
	dll, err := syscall.LoadDLL(path)
	if err != nil {
		return RuntimeVersion{}, err
	}
	defer dll.Release()
	proc, err := dll.FindProc("CEFLibVersion")
	if err != nil {
		return RuntimeVersion{}, fmt.Errorf("CEFLibVersion symbol not found in %s: %w", path, err)
	}
	var fn func(purego.CDecl, *int32, *int32, *int32, *int32)
	purego.RegisterFunc(&fn, proc.Addr())
	var version RuntimeVersion
	fn(purego.CDecl{}, &version.Major, &version.Minor, &version.Release, &version.Build)
	return version, nil
}

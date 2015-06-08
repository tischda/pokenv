// +build windows

package main

import (
	"syscall"
	"unsafe"
)

// Code Copyright 2015 The Go Authors extracted from:
// https://github.com/golang/sys/blob/master/windows/registry/zsyscall_windows.go

var (
	// Advanced Services (advapi32.dll) provide access to the Windows registry
	modadvapi32         = syscall.NewLazyDLL("advapi32.dll")
	procRegSetValueExW  = modadvapi32.NewProc("RegSetValueExW")
	procRegDeleteValueW = modadvapi32.NewProc("RegDeleteValueW")
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724923(v=vs.85).aspx
func regSetValueEx(key syscall.Handle, valueName *uint16, reserved uint32, vtype uint32, buf *byte, bufsize uint32) (regerrno error) {
	r0, _, _ := syscall.Syscall6(procRegSetValueExW.Addr(), 6, uintptr(key), uintptr(unsafe.Pointer(valueName)), uintptr(reserved), uintptr(vtype), uintptr(unsafe.Pointer(buf)), uintptr(bufsize))
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724851(v=vs.85).aspx
func regDeleteValue(key syscall.Handle, name *uint16) (regerrno error) {
	r0, _, _ := syscall.Syscall(procRegDeleteValueW.Addr(), 2, uintptr(key), uintptr(unsafe.Pointer(name)), 0)
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

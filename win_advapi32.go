// +build windows

package main

import (
	"syscall"
	"unsafe"
)

// Code inspired from:
// https://github.com/golang/sys/blob/master/windows/registry/zsyscall_windows.go

var (
	// Advanced Services (advapi32.dll) provide access to the Windows registry
	modadvapi32         = syscall.NewLazyDLL("advapi32.dll")
	procRegSetValueExW  = modadvapi32.NewProc("RegSetValueExW")
	procRegDeleteValueW = modadvapi32.NewProc("RegDeleteValueW")
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724923(v=vs.85).aspx
func regSetValueEx(hKey syscall.Handle, lpValueName *uint16, Reserved uint32, dwType uint32, lpData *byte, cbData uint32) (regerrno error) {
	ret, _, _ := procRegSetValueExW.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(lpValueName)),
		uintptr(Reserved),
		uintptr(dwType),
		uintptr(unsafe.Pointer(lpData)),
		uintptr(cbData))

	// If the function fails, the return value is a nonzero error code defined in Winerror.h
	if ret != 0 {
		regerrno = syscall.Errno(ret)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724851(v=vs.85).aspx
func regDeleteValue(hKey syscall.Handle, lpValueName *uint16) (regerrno error) {
	ret, _, _ := procRegDeleteValueW.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(lpValueName)))

	// If the function fails, the return value is a nonzero error code defined in Winerror.h
	if ret != 0 {
		regerrno = syscall.Errno(ret)
	}
	return
}

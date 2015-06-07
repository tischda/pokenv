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
	modadvapi32 = syscall.NewLazyDLL("advapi32.dll")

	procRegCreateKeyExW = modadvapi32.NewProc("RegCreateKeyExW")
	procRegDeleteKeyW   = modadvapi32.NewProc("RegDeleteKeyW")
	procRegSetValueExW  = modadvapi32.NewProc("RegSetValueExW")
	procRegEnumValueW   = modadvapi32.NewProc("RegEnumValueW")
	procRegDeleteValueW = modadvapi32.NewProc("RegDeleteValueW")
)

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724844(v=vs.85).aspx
func regCreateKeyEx(key syscall.Handle, subkey *uint16, reserved uint32, class *uint16, options uint32, desired uint32, sa *syscall.SecurityAttributes, result *syscall.Handle, disposition *uint32) (regerrno error) {
	r0, _, _ := syscall.Syscall9(procRegCreateKeyExW.Addr(), 9, uintptr(key), uintptr(unsafe.Pointer(subkey)), uintptr(reserved), uintptr(unsafe.Pointer(class)), uintptr(options), uintptr(desired), uintptr(unsafe.Pointer(sa)), uintptr(unsafe.Pointer(result)), uintptr(unsafe.Pointer(disposition)))
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724845(v=vs.85).aspx
func regDeleteKey(key syscall.Handle, subkey *uint16) (regerrno error) {
	r0, _, _ := syscall.Syscall(procRegDeleteKeyW.Addr(), 2, uintptr(key), uintptr(unsafe.Pointer(subkey)), 0)
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

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

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724865(v=vs.85).aspx
func regEnumValue(key syscall.Handle, index uint32, name *uint16, nameLen *uint32, reserved *uint32, valtype *uint32, buf *byte, buflen *uint32) (regerrno error) {
	r0, _, _ := syscall.Syscall9(procRegEnumValueW.Addr(), 8, uintptr(key), uintptr(index), uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(nameLen)), uintptr(unsafe.Pointer(reserved)), uintptr(unsafe.Pointer(valtype)), uintptr(unsafe.Pointer(buf)), uintptr(unsafe.Pointer(buflen)), 0)
	if r0 != 0 {
		regerrno = syscall.Errno(r0)
	}
	return
}

// +build windows

package main

import (
	"log"
	"syscall"
	"unsafe"
)

type realRegistry struct{}

// Writes a REG_QWORD (uint64) to the Windows registry.
func (realRegistry) SetQword(path string, valueName string, value uint64) error {
	handle := OpenKey(path, syscall.KEY_SET_VALUE)
	defer syscall.RegCloseKey(handle)

	return regSetValueEx(
		handle,
		syscall.StringToUTF16Ptr(valueName),
		0,
		syscall.REG_QWORD,
		(*byte)(unsafe.Pointer(&value)),
		8)
}

// Opens a Windows HKCU registry key and returns a handle. You must close
// the handle with `defer syscall.RegCloseKey(handle)` in the calling code.
func OpenKey(path string, desiredAccess uint32) syscall.Handle {
	var handle syscall.Handle
	err := syscall.RegOpenKeyEx(
		syscall.HKEY_CURRENT_USER,
		syscall.StringToUTF16Ptr(path),
		0,
		desiredAccess,
		&handle)
	if err != nil {
		log.Printf("Cannot open path %q in the registry\n", path)
	}
	return handle
}

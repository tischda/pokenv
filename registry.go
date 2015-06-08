// +build windows

package main

import (
	"log"
	"strings"
	"syscall"
	"unsafe"
)

type realRegistry struct{}

type regPath struct {
	key    syscall.Handle
	subKey string
}

// List of authorized keys, indexed by enum constant in main.go
var target []regPath = []regPath{
	regPath{syscall.HKEY_CURRENT_USER, `Environment`},
	regPath{syscall.HKEY_LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`},
}

// Writes a string to the Windows registry. Type is set to REG_EXPAND_SZ when
// the value contains "%", otherwise it will use REG_SZ.
func (realRegistry) SetString(key int, valueName string, value string) error {
	handle := OpenKey(key, syscall.KEY_SET_VALUE)
	defer syscall.RegCloseKey(handle)

	// set type
	var valueType uint32 = syscall.REG_SZ
	if strings.Contains(value, "%") {
		valueType = syscall.REG_EXPAND_SZ
	}

	return regSetValueEx(
		handle,
		syscall.StringToUTF16Ptr(valueName),
		0,
		valueType,
		(*byte)(unsafe.Pointer(syscall.StringToUTF16Ptr(value))),
		uint32(len(value)*2))
}

// Deletes a key value from the Windows registry.
func (realRegistry) DeleteValue(key int, valueName string) error {
	handle := OpenKey(key, syscall.KEY_WRITE)
	defer syscall.RegCloseKey(handle)
	return regDeleteValue(handle, syscall.StringToUTF16Ptr(valueName))
}

// Opens a Windows HKCU or HKLM registry key and returns a handle. You must close
// the handle with `defer syscall.RegCloseKey(handle)` in the calling code.
func OpenKey(key int, desiredAccess uint32) syscall.Handle {
	var handle syscall.Handle
	err := syscall.RegOpenKeyEx(
		target[key].key,
		syscall.StringToUTF16Ptr(target[key].subKey),
		0,
		desiredAccess,
		&handle)
	if err != nil {
		log.Printf("Cannot open path %q in the registry\n", target[key])
	}
	return handle
}

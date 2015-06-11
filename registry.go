// +build windows

package main

import (
	"log"
	"strings"
	"syscall"
	"unsafe"
)

type realRegistry struct{}

// do not reorder
var hKeyTable = []syscall.Handle{
	syscall.HKEY_CLASSES_ROOT,
	syscall.HKEY_CURRENT_USER,
	syscall.HKEY_LOCAL_MACHINE,
	syscall.HKEY_USERS,
	syscall.HKEY_PERFORMANCE_DATA,
	syscall.HKEY_CURRENT_CONFIG,
	syscall.HKEY_DYN_DATA,
}

// Writes a string to the Windows registry. Type is set to REG_EXPAND_SZ when
// the value contains "%", otherwise it will use REG_SZ.
func (realRegistry) SetString(path regPath, valueName string, value string) error {
	handle := openKey(path, syscall.KEY_SET_VALUE)
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
func (realRegistry) DeleteValue(path regPath, valueName string) error {
	handle := openKey(path, syscall.KEY_WRITE)
	defer syscall.RegCloseKey(handle)

	return regDeleteValue(handle, syscall.StringToUTF16Ptr(valueName))
}

// Opens a Windows registry key and returns a handle. You must close
// the handle with `defer syscall.RegCloseKey(handle)` in the calling code.
func openKey(path regPath, desiredAccess uint32) syscall.Handle {
	var handle syscall.Handle

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724897(v=vs.85).aspx
	err := syscall.RegOpenKeyEx(
		hKeyTable[path.hKeyIdx],
		syscall.StringToUTF16Ptr(path.lpSubKey),
		0,
		desiredAccess,
		&handle)

	if err != nil {
		log.Fatalln("Cannot open registry path:", path)
	}
	return handle
}

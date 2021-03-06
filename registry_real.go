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
func (realRegistry) SetString(path regKey, valueName string, value string) error {
	handle := openKey(path, syscall.KEY_SET_VALUE)
	defer syscall.RegCloseKey(handle)

	// set type
	var valueType uint32 = syscall.REG_SZ
	if strings.Contains(value, "%") {
		valueType = syscall.REG_EXPAND_SZ
	}

	// set value
	return regSetValueEx(
		handle,
		StringToUTF16Ptr(valueName),
		0,
		valueType,
		(*byte)(unsafe.Pointer(StringToUTF16Ptr(value))),
		uint32(len(value)*2))
}

// Read string from Windows registry (no expansion).
// Thanks to http://npf.io/2012/11/go-win-stuff/
func (realRegistry) GetString(path regKey, valueName string) (value string, err error) {
	handle := openKey(path, syscall.KEY_QUERY_VALUE)
	defer syscall.RegCloseKey(handle)

	var typ uint32
	var bufSize uint32

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724911(v=vs.85).aspx
	err = syscall.RegQueryValueEx(
		handle,
		StringToUTF16Ptr(valueName),
		nil,
		&typ,
		nil,
		&bufSize)

	if err != nil {
		return
	}

	data := make([]uint16, bufSize/2+1)

	err = syscall.RegQueryValueEx(
		handle,
		StringToUTF16Ptr(valueName),
		nil,
		&typ,
		(*byte)(unsafe.Pointer(&data[0])),
		&bufSize)

	if err != nil {
		return
	}
	return syscall.UTF16ToString(data), nil
}

// Deletes a key value from the Windows registry.
func (realRegistry) DeleteValue(path regKey, valueName string) error {
	handle := openKey(path, syscall.KEY_SET_VALUE)
	defer syscall.RegCloseKey(handle)

	return regDeleteValue(handle, StringToUTF16Ptr(valueName))
}

// Opens a Windows registry key and returns a handle. You must close
// the handle with `defer syscall.RegCloseKey(handle)` in the calling code.
func openKey(path regKey, desiredAccess uint32) syscall.Handle {
	var handle syscall.Handle

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms724897(v=vs.85).aspx
	err := syscall.RegOpenKeyEx(
		hKeyTable[path.hKeyIdx],
		StringToUTF16Ptr(path.lpSubKey),
		0,
		desiredAccess,
		&handle)

	if err != nil {
		log.Fatalln("Cannot open registry path:", path)
	}
	return handle
}

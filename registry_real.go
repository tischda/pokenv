//go:build windows

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
	defer func() {
		err := syscall.RegCloseKey(handle)
		if err != nil {
			log.Printf("pokenv: failed to close registry key: %v", err)
		}
	}()

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

// Deletes a key value from the Windows registry.
func (realRegistry) DeleteValue(path regKey, valueName string) error {
	handle := openKey(path, syscall.KEY_SET_VALUE)
	defer func() {
		err := syscall.RegCloseKey(handle)
		if err != nil {
			log.Printf("pokenv: failed to close registry key: %v", err)
		}
	}()

	return regDeleteValue(handle, StringToUTF16Ptr(valueName))
}

// Opens a Windows registry key and returns a handle. You must close the
// handle with `defer syscall.RegCloseKey(handle)` in the calling code.
func openKey(path regKey, desiredAccess uint32) syscall.Handle {
	var handle syscall.Handle

	subkey, err := syscall.UTF16PtrFromString(path.lpSubKey)
	if err != nil {
		log.Fatalln("Error on registry path.subKey:", path.lpSubKey, err)
	}

	err = syscall.RegOpenKeyEx(
		hKeyTable[path.hKeyIdx],
		subkey,
		0,
		desiredAccess,
		&handle)

	if err != nil {
		log.Fatalln("Cannot open registry path:", path, err)
	}
	return handle
}

// refresh Environment
func refresh() error {
	lParam, _ := syscall.UTF16PtrFromString("Environment")

	// note that when an application sends this message, wParam must be NULL:
	// https://learn.microsoft.com/en-us/windows/win32/winmsg/wm-settingchange
	ret := SendMessageTimeout(HWND_BROADCAST, WM_SETTINGCHANGE, nil,
		lParam, SMTO_NORMAL|SMTO_ABORTIFHUNG, TIMEOUT_MS)

	if ret == 0 { // if the function succeeds, the return value is non-zero
		return syscall.GetLastError()
	}
	return nil
}

// https://golang.org/src/syscall/syscall_windows.go
// syscall.StringToUTF16Ptr is deprecated, here is our own:
func StringToUTF16Ptr(s string) *uint16 {
	ptr, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		log.Fatalln("String with NULL passed to StringToUTF16Ptr")
	}
	return ptr
}

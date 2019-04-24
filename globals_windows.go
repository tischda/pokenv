// +build windows

package main

import "unsafe"

var registry = realRegistry{}

func refreshEnvironment() {
	var ptr = unsafe.Pointer(StringToUTF16Ptr("Environment"))
	SendMessageTimeout(HWND_BROADCAST, WM_SETTINGCHANGE, 0, uintptr(ptr), SMTO_ABORTIFHUNG, 5000)
}

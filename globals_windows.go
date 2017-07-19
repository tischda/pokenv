// +build windows

package main

import "unsafe"

var registry = realRegistry{}

func refreshEnvironment() {
	SendMessageTimeout(HWND_BROADCAST, WM_SETTINGCHANGE, 0,
		uintptr(unsafe.Pointer(StringToUTF16Ptr("Environment"))), SMTO_ABORTIFHUNG, 5000)
}

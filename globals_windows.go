//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"
)

var registry = realRegistry{}

// When an application sends this message, wParam must be NULL:
// https://learn.microsoft.com/en-us/windows/win32/winmsg/wm-settingchange
func refreshEnvironment() {
	ret := SendMessageTimeout(HWND_BROADCAST, WM_SETTINGCHANGE, nil,
		StringToUTF16Ptr("Environment"), SMTO_NORMAL|SMTO_ABORTIFHUNG, TIMEOUT_MS)

	// If the function succeeds, the return value is nonzero
	if ret == 0 {
		fmt.Println("Refresh: Error")

		//TODO: don't exit here, but return error
		os.Exit(1)
	}
}

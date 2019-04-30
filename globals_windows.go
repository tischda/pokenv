// +build windows

package main

import (
	"log"
)

var registry = realRegistry{}

// timeout in milliseconds
func refreshEnvironment() {
	ret := SendMessageTimeout(HWND_BROADCAST, WM_SETTINGCHANGE, StringToUTF16Ptr(""),
		StringToUTF16Ptr("Environment"), SMTO_NORMAL|SMTO_ABORTIFHUNG, 5000)

	// If the function succeeds, the return value is nonzero
	if ret == 0 {
		log.Fatalln("Refresh: Error")
	}
}

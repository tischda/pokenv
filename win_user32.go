// +build windows

package main

import (
	"log"
	"syscall"
)

var (
	moduser32 = syscall.NewLazyDLL("user32.dll")

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms644952(v=vs.85).aspx
	procSendMessageTimeout = moduser32.NewProc("SendMessageTimeoutW")
)

type HWND uintptr

const (
	HWND_BROADCAST = HWND(0xFFFF)

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms725497(v=vs.85).aspx
	WM_WININICHANGE  = 0x001A
	WM_SETTINGCHANGE = WM_WININICHANGE

	// https://msdn.microsoft.com/en-us/library/windows/desktop/ms644952(v=vs.85).aspx
	SMTO_ABORTIFHUNG = 0x0002
)

// https://github.com/AllenDang/w32/blob/master/user32.go#L318
func SendMessageTimeout(hwnd HWND, msg uint32, wParam, lParam uintptr, fuFlags, uTimeout uint32) uintptr {
	ret, _, _ := procSendMessageTimeout.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam,
		uintptr(fuFlags),
		uintptr(uTimeout))

	return ret
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

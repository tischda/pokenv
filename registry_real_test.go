//go:build windows

package main

import (
	"syscall"
	"testing"
	"unsafe"
)

func TestSetDeleteValue(t *testing.T) {

	var registry = realRegistry{}

	expected := "titi"

	// set value
	err := registry.SetString(REG_KEY_USER, "toto", expected)
	if err != nil {
		t.Errorf("Error in SetString: %q", err)
	}

	// read value
	actual, err := getString(REG_KEY_USER, "toto")
	if err != nil {
		t.Errorf("Error in GetString: %q", err)
	}
	if actual != expected {
		t.Errorf("Expected: %q, was: %q", expected, actual)
	}

	// delete value
	err = registry.DeleteValue(REG_KEY_USER, "toto")
	if err != nil {
		t.Errorf("Error in DeleteValue: %q", err)
	}
}

// Read string from Windows registry (no expansion).
func getString(path regKey, valueName string) (value string, err error) {
	handle := openKey(path, syscall.KEY_QUERY_VALUE)
	defer syscall.RegCloseKey(handle)

	var typ uint32
	var bufSize uint32

	name, err := syscall.UTF16PtrFromString(valueName)
	if err != nil {
		return "", err
	}

	// First call: Get the required buffer size
	// Pass nil for data buffer to get size in bufSize
	err = syscall.RegQueryValueEx(
		handle,
		name,
		nil,
		&typ,
		nil,      // nil data buffer
		&bufSize) // receives required size

	if err != nil {
		return "", err
	}

	// Allocate buffer with the exact size needed
	// Add 1 to handle potential rounding for UTF16
	data := make([]uint16, bufSize/2+1)

	// Second call: Actually get the data with properly sized buffer
	err = syscall.RegQueryValueEx(
		handle,
		name,
		nil,
		&typ,
		(*byte)(unsafe.Pointer(&data[0])), // properly sized buffer
		&bufSize)

	if err != nil {
		return "", err
	}
	return syscall.UTF16ToString(data), nil
}

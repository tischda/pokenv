package main

// A registry path is composed of an hKey index and the string representation
// of the path withing that hKey. We use hKey indexes to avoid dependency on
// non-portable syscall values.
type regKey struct {
	hKeyIdx  uint8
	lpSubKey string
}

// Registry hKey index values, do not reorder
const (
	HKEY_CLASSES_ROOT = iota
	HKEY_CURRENT_USER
	HKEY_LOCAL_MACHINE
	HKEY_USERS
	HKEY_PERFORMANCE_DATA
	HKEY_CURRENT_CONFIG
	HKEY_DYN_DATA
)

type Registry interface {
	DeleteValue(path regKey, valueName string) error
	SetString(path regKey, valueName string, value string) error
}

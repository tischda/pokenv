package main

// Registry key indexes, do not reorder
const (
	HKCU = iota
	HKLM
)

type Registry interface {
	DeleteValue(key int, valueName string) error
	SetString(key int, valueName string, value string) error
}

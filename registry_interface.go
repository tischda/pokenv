package main

type Registry interface {
	DeleteValue(key int, valueName string) error
	SetString(key int, valueName string, value string) error
}

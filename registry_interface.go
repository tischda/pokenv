package main

type Registry interface {
	SetString(key int, valueName string, value string) error
}

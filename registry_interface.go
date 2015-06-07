package main

type Registry interface {
	SetString(path string, valueName string, value string) error
}

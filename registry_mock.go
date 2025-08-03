package main

import (
	"errors"
)

type mockRegistry struct {
	env map[string]string
}

var mock = mockRegistry{}

func init() {
	mock.env = make(map[string]string)
}

func (r mockRegistry) SetString(path regKey, valueName string, value string) error {
	r.env[valueName] = value
	return nil
}

//lint:file-ignore ST1005 : this is the original Windows error message
func (r mockRegistry) DeleteValue(path regKey, valueName string) error {
	if _, exists := r.env[valueName]; !exists {
		return errors.New("The system cannot find the file specified.")
	}
	delete(r.env, valueName)
	return nil
}

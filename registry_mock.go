package main

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

func (r mockRegistry) DeleteValue(path regKey, valueName string) error {
	delete(r.env, valueName)
	return nil
}

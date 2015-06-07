package main

type mockRegistry struct {
	env map[string]string
}

func (r mockRegistry) SetString(path string, valueName string, value string) error {
	r.env[valueName] = value
	return nil
}

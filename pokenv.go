package main

import (
	"log"
	"os"
	"regexp"
	"strings"
)

var REG_KEY_USER = regKey{HKEY_CURRENT_USER, `Environment`}

// TODO: access to this key requires admin rights (cf. https://github.com/mozey/run-as-admin)
var REG_KEY_MACHINE = regKey{HKEY_LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`}

type pokenv struct {
	registry  Registry
	checkPath bool
}

func (p *pokenv) processFile(reg regKey, fileName string) error {
	vars := p.parseFile(fileName)

	if !p.checkPath || assertValuesAreValidPaths(&vars) {
		return p.setVars(reg, vars)
	}
	return nil
}

func (p *pokenv) parseFile(fileName string) varMap {
	var file *os.File

	if fileName == "stdin" {
		file = os.Stdin
	} else {
		var err error
		file, err = os.Open(fileName) // O_RDONLY mode
		if err != nil {
			log.Fatalln(err)
		}
	}
	defer file.Close()

	var parser = &parser{}
	return parser.parse(file)
}

func (p *pokenv) setVars(reg regKey, vars varMap) error {
	for variable, values := range vars {
		var err error
		if strings.HasPrefix(variable, "-") {
			log.Println("Deleting", variable[1:])
			err = p.registry.DeleteValue(reg, variable[1:])
			// if the variable does not exist, we ignore the error
			if err != nil {
				log.Printf("Warning: %s\n", err)
				err = nil // continue processing
			}
		} else if len(values) == 0 {
			log.Printf("Warning: [%s] is empty, will be left untouched\n", variable)
		} else {
			value := strings.Join(values, ";")
			log.Printf("Setting `%s` to `%s`\n", variable, value)
			err = p.registry.SetString(reg, variable, value)
		}
		if err != nil {
			return (err) // leave loop on first error
		}
	}
	return nil
}

// checks if path is valid.
// Does Windows variable expansion so that '%windir%' resolves to 'c:\Windows'.
func isPathInvalid(path string) bool {
	for strings.Contains(path, "%") {
		exp := regexp.MustCompile(`(.*)%(.*)%(.*)`)
		parts := exp.FindStringSubmatch(path)
		path = parts[1] + os.ExpandEnv("${"+parts[2]+"}") + parts[3]
	}
	_, err := os.Stat(path)
	return err != nil
}

// checks all values in the vars map assuming they are valid paths (all values are checked)
// returns:
// * false if at least one path is invalid
// * true if all values are valid paths
func assertValuesAreValidPaths(vars *varMap) bool {
	ret := true
	for key, value := range *vars {
		for _, line := range value {
			if isPathInvalid(line) {
				log.Printf("Invalid path in section [%s]: %s\n", key, line)
				ret = false
			}
		}
	}
	return ret
}

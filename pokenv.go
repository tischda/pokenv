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

func (p *pokenv) processFile(reg regKey, fileName string) {
	vars := p.parseFile(fileName)

	if !p.checkPath || assertValuesAreValidPaths(&vars) {
		p.setVars(reg, vars)
	}
}

func (p *pokenv) parseFile(fileName string) varMap {
	var file *os.File

	if fileName == "REQUIRED" {
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

// TODO: log.Fatalln() exits, this is too strong! It's probably better just to
// log the issue and let the function return an error after processing the other
// variables.
func (p *pokenv) setVars(reg regKey, vars varMap) {
	for variable, values := range vars {
		if len(values) == 0 {
			log.Println("Deleting", variable)
			err := p.registry.DeleteValue(reg, variable)
			if err != nil {
				log.Fatalln(err)
			}

		} else {
			joined := strings.Join(values, ";")
			log.Printf("Setting `%s` to `%s`\n", variable, joined)
			err := p.registry.SetString(reg, variable, joined)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
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

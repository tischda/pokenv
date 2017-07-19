package main

import (
	"fmt"
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

func (p *pokenv) processFile(reg regKey, file *os.File) {
	parser := &parser{}
	vars := parser.parse(file)
	ok := true
	if p.checkPath {
		ok = assertValuesAreValidPaths(&vars)
	}
	if ok {
		p.setVars(reg, vars)
	}
}

func (p *pokenv) setVars(reg regKey, vars varMap) {
	for variable, values := range vars {
		if len(values) == 0 {
			log.Println("Deleting", variable)
			p.registry.DeleteValue(reg, variable)
		} else {
			joined := strings.Join(values, ";")
			log.Printf("Setting `%s` to `%s`\n", variable, joined)
			p.registry.SetString(reg, variable, joined)
		}
	}
}

// checks if path is valid.
// Does Windows variable expansion so that '%windir%' resolves to 'c:\Windows'.
func isPathInvalid(path string) bool {
	for strings.Contains(path, "%") {
		regexp := regexp.MustCompile(`(.*)%(.*)%(.*)`)
		parts := regexp.FindStringSubmatch(path)
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
				log.Print(fmt.Sprintf("Invalid path in section [%s]: %s\n", key, line))
				ret = false
			}
		}
	}
	return ret
}

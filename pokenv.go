package main

import (
	"log"
	"os"
	"regexp"
	"strings"
)

var PATH_USER = regPath{HKEY_CURRENT_USER, `Environment`}
var PATH_MACHINE = regPath{HKEY_LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`}

type pokenv struct {
	registry  Registry
	pathcheck bool
}

func (p *pokenv) importFromFile(path regPath, fileName string) {
	vars := p.processFile(fileName)
	if p.pathcheck {
		for _, path := range values(vars) {
			if isPathInvalid(path) {
				log.Fatalln("Invalid path:", path)
			}
		}
	}
	p.setVars(path, vars)
}

func (p *pokenv) processFile(fileName string) varMap {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	parser := &parser{}
	return parser.processAllLines(file)
}

func (p *pokenv) setVars(path regPath, vars varMap) {
	for variable, values := range vars {
		if len(values) == 0 {
			log.Println("Deleting", variable)
			p.registry.DeleteValue(path, variable)
		} else {
			joined := strings.Join(values, ";")
			log.Printf("Setting `%s` to `%s`\n", variable, joined)
			p.registry.SetString(path, variable, joined)
		}
	}
}

// checks if path is valid. Does Windows variable expansion
// so that '%windir%' resolves to 'c:\Windows'.
func isPathInvalid(value string) bool {
	var filename = value
	for strings.Contains(filename, "%") {
		regexp := regexp.MustCompile(`(.*)%(.*)%(.*)`)
		parts := regexp.FindStringSubmatch(filename)
		filename = parts[1] + os.ExpandEnv("${"+parts[2]+"}") + parts[3]
	}
	_, err := os.Stat(filename)
	return err != nil
}

func values(m varMap) []string {
	v := make([]string, 1)

	for _, lines := range m {
		v = append(v, lines...)
	}
	return v
}

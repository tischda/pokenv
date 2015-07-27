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
	env := p.processFile(fileName)
	// validate paths if pathcheck
	p.setVars(path, env)
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

func (p *pokenv) setVars(path regPath, env varMap) {
	for variable, values := range env {
		if firstValueIsEmpty(values) {
			log.Println("Deleting", variable)
			p.registry.DeleteValue(path, variable)
		} else {
			joined := strings.Join(values, ";")
			log.Printf("Setting `%s` to `%s`\n", variable, joined)
			p.registry.SetString(path, variable, joined)
		}
	}
}

//if p.pathcheck {
//	if isPathInvalid(value) {
//		log.Fatalln("Invalid path:", value)
//	}
//}

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

func firstValueIsEmpty(values []string) bool {
	return values[0] == ""
}

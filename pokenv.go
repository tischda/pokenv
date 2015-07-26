package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
	"io"
)

var sectionRegex = regexp.MustCompile(`^\[(.*)\]$`)

var PATH_USER = regPath{HKEY_CURRENT_USER, `Environment`}
var PATH_MACHINE = regPath{HKEY_LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`}

type pokenv struct {
	environment      map[string][]string
	setContainsValue map[string]bool
	currentVariable  string
	registry         Registry
	pathcheck        bool
}

func (p *pokenv) importEnv(path regPath, fileName string) {
	p.processFile(fileName)
	p.setVars(path)
}

func (p *pokenv) processFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	p.processAllLines(file)
}

func (p *pokenv) setVars(path regPath) {
	for variable, values := range p.environment {
		if p.firstValueIsEmpty(values) {
			log.Println("Deleting", variable)
			p.registry.DeleteValue(path, variable)
		} else {
			joined := strings.Join(values, ";")
			log.Printf("Setting `%s` to `%s`\n", variable, joined)
			p.registry.SetString(path, variable, joined)
		}
	}
}

func (p *pokenv) processAllLines(r io.Reader) {
	p.environment = make(map[string][]string)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		p.processLine(strings.TrimSpace(line))
	}
}

func (p *pokenv) processLine(line string) {
	// ignore blank line
	if line == "" {
		return
	}
	// ignore comment
	if strings.HasPrefix(line, "#") {
		return
	}
	// sections and values
	match := sectionRegex.FindStringSubmatch(line)
	if match != nil {
		p.processSection(match[1])
	} else {
		p.processValue(trimComments(line))
	}
}

func (p *pokenv) processSection(section string) {
	p.currentVariable = strings.Replace(section, " ", "", -1)

	// mark section as empty (required for deletion)
	if p.environment[p.currentVariable] == nil {
		p.addToCurrentVariable("")
	}
}

func (p *pokenv) processValue(value string) {
	if p.setContainsValue[value] {
		log.Println("Warning: duplicate entry:", value)
	} else {
		if p.pathcheck {
			if isPathInvalid(value) {
				log.Fatalln("Invalid path:", value)
			}
		}

		// if this is first value, initialize and set value
		if p.firstValueIsEmpty(p.environment[p.currentVariable]) {
			p.setContainsValue = make(map[string]bool)
			p.setFirst(value)
		} else {
			p.addToCurrentVariable(value)
		}
		p.setContainsValue[value] = true
	}
}

func (p *pokenv) firstValueIsEmpty(values []string) bool {
	return values[0] == ""
}

func (p *pokenv) setFirst(value string) {
	p.environment[p.currentVariable][0] = value
}

// add value to current variable in environment
func (p *pokenv) addToCurrentVariable(value string) {
	p.environment[p.currentVariable] = append(p.environment[p.currentVariable], value)
}

func trimComments(s string) string {
	if idx := strings.Index(s, "#"); idx != -1 {
		return strings.TrimSpace(s[:idx])
	}
	return s
}

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

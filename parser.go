package main

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"
)

var isSectionRegex = regexp.MustCompile(`^\[(.*)\]$`)

type varMap map[string][]string
type stringSet map[string]bool

type parser struct {
	vars       varMap
	currentVar string
	currentSet stringSet
}

func (p *parser) parse(r io.Reader) varMap {
	p.cleanUp()
	p.vars = make(varMap)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		p.parseLine(strings.TrimSpace(line))
	}
	p.closePreviousSectionIfEmpty()
	p.cleanUp()
	return p.vars
}

func (p *parser) parseLine(line string) {
	// ignore blank line
	if line == "" {
		return
	}
	// ignore comment
	if strings.HasPrefix(line, "#") {
		return
	}
	// sections and values
	match := isSectionRegex.FindStringSubmatch(line)
	if match != nil {
		p.parseSection(match[1])
	} else {
		p.parseValue(line)
	}
}

func (p *parser) parseSection(section string) {
	// this is new section, close previous
	p.closePreviousSectionIfEmpty()

	// start clean
	p.currentVar = trimSpaces(section)
	p.currentSet = make(stringSet)

	// if section exists, restore duplicates list
	values := p.vars[p.currentVar]
	if len(values) > 0 {
		for _, v := range values {
			p.currentSet[v] = true
		}
	}
}

func (p *parser) parseValue(value string) {
	if p.currentVar == "" {
		log.Println("Error: orphan line (not in section):", value)
	} else {
		value = trimComments(value)
		if p.currentSet[value] {
			log.Println("Warning: duplicate entry:", value)
		} else {
			p.vars[p.currentVar] = append(p.vars[p.currentVar], value)
			p.currentSet[value] = true
		}
	}
}

// If section is empty, add empty list to mark for deletion.
func (p *parser) closePreviousSectionIfEmpty() {
	if p.currentVar != "" && len(p.currentSet) == 0 {
		p.vars[p.currentVar] = []string{}
	}
}

func (p *parser) cleanUp() {
	p.currentVar = ""
	p.currentSet = nil
}

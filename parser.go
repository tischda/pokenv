package main

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"
)

var isSectionRegex = regexp.MustCompile(`^\[(.*)\]$`)

// TODO: reused in pokenv
type varMap map[string][]string

type stringSet map[string]bool

type parser struct {
	vars       varMap
	currentVar string
	currentSet stringSet
}

func (p *parser) processAllLines(r io.Reader) varMap {
	p.vars = make(varMap)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		p.processLine(strings.TrimSpace(line))
	}
	p.setVars()
	return p.vars
}

func (p *parser) processLine(line string) {
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
		p.processSection(match[1])
	} else {
		p.processValue(line)
	}
}

func (p *parser) processSection(section string) {

	// new section, dump recorded contents to vars
	p.setVars()

	// start clean
	p.currentVar = trimSpaces(section)
	p.currentSet = make(stringSet)

	// already had that section ?
	values := p.vars[p.currentVar]
	if values != nil {
		for _, v := range values {
			p.currentSet[v] = true
		}
	}
}

func (p *parser) processValue(value string) {
	if p.currentVar == "" {
		log.Println("Orphan line (not in section):", value)
	} else {
		value = trimComments(value)
		if p.currentSet[value] {
			log.Println("Warning: duplicate entry:", value)
		} else {
			p.currentSet[value] = true
		}
	}
}

// Copies the parsed variable to the parser's varMap.
func (p *parser) setVars() {
	if p.currentVar != "" {
		p.vars[p.currentVar] = keys(p.currentSet)
		p.currentVar = ""
		p.currentSet = nil
	}
}

// ------------------ utility

// If there is a comment in the line, return only
// content to the left of the comment marker '#'.
func trimComments(s string) string {
	if idx := strings.Index(s, "#"); idx != -1 {
		return strings.TrimSpace(s[:idx])
	}
	return s
}

// Remove all spaces in line.
func trimSpaces(s string) string {
	return strings.Replace(s, " ", "", -1)
}

func keys(set stringSet) []string {
	keys := make([]string, len(set))
	i := 0
	for k := range set {
		keys[i] = k
		i += 1
	}
	return keys
}

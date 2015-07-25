// +build windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// http://technosophos.com/2014/06/11/compile-time-string-in-go.html
// go build -ldflags "-x main.version $(git describe --tags)"
var version string

func main() {
	hkcu := flag.String("hkcu", "REQUIRED", "process input file into HKEY_CURRENT_USER environment")
	hklm := flag.String("hklm", "REQUIRED", "process input file into HKEY_LOCAL_MACHINE environment")
	check := flag.Bool("checkpaths", false, "values are paths, check if they are valid on this system")
	showVersion := flag.Bool("version", false, "print version and exit")

	// configure logging
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-checkpaths] [-hkcu|-hklm] infile\n  infile: the input file\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Println("pokenv version", version)
		return
	}
	if flag.NFlag() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	p := pokenv{
		environment: make(map[string][]string),
		registry:    realRegistry{},
		pathcheck:   *check,
	}

	if *hkcu != "REQUIRED" {
		p.importEnv(PATH_USER, *hkcu)
	}
	if *hklm != "REQUIRED" {
		p.importEnv(PATH_MACHINE, *hklm)
	}
}

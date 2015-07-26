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

var hkcu string
var hklm string
var check bool
var showVersion bool

func init() {
	flag.StringVar(&hkcu, "hkcu", "REQUIRED", "process input file into HKEY_CURRENT_USER environment")
	flag.StringVar(&hklm, "hklm", "REQUIRED", "process input file into HKEY_LOCAL_MACHINE environment")
	flag.BoolVar(&check, "checkpaths", false, "values are paths, check if they are valid on this system")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
}

func main() {
	// configure logging
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-checkpaths] [-hkcu|-hklm] infile\n  infile: the input file\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if showVersion {
		fmt.Println("pokenv version", version)
	} else {
		if flag.NFlag() < 1 {
			flag.Usage()
			os.Exit(1)
		}

		// p is defined in global_xxx.go

		if hkcu != "REQUIRED" {
			p.importEnv(PATH_USER, hkcu)
		}
		if hklm != "REQUIRED" {
			p.importEnv(PATH_MACHINE, hklm)
		}
	}
}

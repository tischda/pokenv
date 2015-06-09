// +build windows

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const version string = "1.0.0"

var registry Registry

const (
	// do not reorder
	HKCU = iota
	HKLM
)

func main() {
	hkcu := flag.Bool("hkcu", false, "set HKEY_CURRENT_USER environment")
	hklm := flag.Bool("hklm", false, "set HKEY_LOCAL_MACHINE environment")
	showVersion := flag.Bool("version", false, "print version and exit")

	registry = realRegistry{}

	// configure logging
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] infile\n  infile: the input file\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Println("pokenv version", version)
		return
	}
	if flag.NArg() != 1 || flag.NFlag() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	if *hkcu {
		setEnv(HKCU, flag.Arg(0))
	}
	if *hklm {
		setEnv(HKLM, flag.Arg(0))
	}
}

// +build windows darwin

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	HKCU string = `HKEY_CURRENT_USER\Environment`
	HKLM string = `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Session Manager\Environment`
)

func main() {
	hkcu := flag.Bool("hkcu", false, "HKEY_CURRENT_USER")
	hklm := flag.Bool("hklm", false, "HKEY_LOCAL_MACHINE")
	fileName := flag.String("f", "REQUIRED", "file name")
	version := flag.Bool("version", false, "print version")

	// registry = realRegistry{}
	registry = mockRegistry{}
	registry.env = make(map[string]string)

	// configure logging
	log.SetFlags(0)

	// parse command line arguments
	flag.Parse()

	if *version {
		fmt.Println("Pokenv v0.1.0")
		return
	}

	if flag.NFlag() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	if *hkcu {
		setEnv(HKCU, *fileName)
	}

	if *hklm {
		setEnv(HKLM, *fileName)
	}
}

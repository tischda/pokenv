package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const PROG_NAME string = "pokenv"

// The duration of the time-out period, in milliseconds. If the message is a broadcast message,
// each window can use the full time-out period:
// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-sendmessagetimeouta
const TIMEOUT_MS = 5000

var version string

var fileName string
var flagHelp = flag.Bool("help", false, "displays this help message")
var flagMachine = flag.Bool("machine", false, "specifies that the variables should be set system wide (HKEY_LOCAL_MACHINE)")
var flagCheck = flag.Bool("checkpaths", false, "check if values are valid paths on this system")
var flagVersion = flag.Bool("version", false, "print version and exit")

func init() {
	flag.BoolVar(flagHelp, "h", false, "")
	flag.BoolVar(flagMachine, "m", false, "")
	flag.BoolVar(flagCheck, "c", false, "")
	flag.BoolVar(flagVersion, "v", false, "")
	flag.StringVar(&fileName, "f", "REQUIRED", "file containing the variables to load into the Windows environment")
}

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-h] [-c] [-m] [-f infile]\n\nOPTIONS:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *flagVersion {
		fmt.Println("pokenv version", version)
	} else {
		if *flagHelp || flag.NArg() > 0 {
			flag.Usage()
			os.Exit(1)
		}
		process()
	}
}

func process() {
	pokenv := pokenv{registry: registry, checkPath: *flagCheck}
	if *flagMachine {
		pokenv.processFile(REG_KEY_MACHINE, fileName)
	} else {
		pokenv.processFile(REG_KEY_USER, fileName)
	}
	refreshEnvironment()
}

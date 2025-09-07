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
var flagCheck = flag.Bool("checkpaths", false, "check if ALL values are valid paths on this system")
var flagVersion = flag.Bool("version", false, "print version and exit")

func init() {
	flag.BoolVar(flagHelp, "h", false, "")
	flag.BoolVar(flagMachine, "m", false, "")
	flag.BoolVar(flagCheck, "c", false, "")
	flag.BoolVar(flagVersion, "v", false, "")
	flag.StringVar(&fileName, "f", "stdin", "file containing the variables to load into the Windows environment")
}

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: "+PROG_NAME+` [--checkpaths] [--machine] [-f inifile] 

where inifile is a file containing the variables to load into the Windows environment.
For example, you can use a file like this:

  [GOBIN]
  c:\usr\bin

  [-GOPATH]
  # prefix with dash '-' to remove variable


OPTIONS:`)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\n\nEXAMPLES:")
		fmt.Fprintln(os.Stderr, "  "+PROG_NAME+` --checkpaths -f paths.ini

  where paths.ini could be the example shown above. When using the --checkpaths option,
  make sure ALL values are actually paths. If you set GOFLAGS here, it will fail.`)
	}
	flag.Parse()

	if flag.Arg(0) == "version" || *flagVersion {
		fmt.Printf("%s version %s\n", PROG_NAME, version)
		return
	}

	if *flagHelp {
		flag.Usage()
		return
	}

	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}

	// process registry
	registry := realRegistry{}
	pokenv := pokenv{registry: registry, checkPath: *flagCheck}
	var err error
	if *flagMachine {
		err = pokenv.processFile(REG_KEY_MACHINE, fileName)
	} else {
		err = pokenv.processFile(REG_KEY_USER, fileName)
	}
	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}

	// refresh Environment
	err = refresh()
	if err != nil {
		log.Println("Refresh error:", err)
		os.Exit(1)
	}
}

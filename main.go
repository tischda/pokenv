package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var version string

var file_name string
var flag_help = flag.Bool("help", false, "displays this help message")
var flag_machine = flag.Bool("machine", false, "specifies that the variables should be set system wide (HKEY_LOCAL_MACHINE)")
var flag_check = flag.Bool("checkpaths", false, "check if values are valid paths on this system")
var flag_version = flag.Bool("version", false, "print version and exit")

func init() {
	flag.BoolVar(flag_help, "h", false, "")
	flag.BoolVar(flag_machine, "m", false, "")
	flag.BoolVar(flag_check, "c", false, "")
	flag.BoolVar(flag_version, "v", false, "")
	flag.StringVar(&file_name, "f", "REQUIRED", "file containing the variables to load into the Windows environment")
}

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-h] [-c] [-m] [-f infile]\n\nOPTIONS:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *flag_version {
		fmt.Println("pokenv version", version)
	} else {
		if *flag_help || flag.NArg() > 0 {
			flag.Usage()
			os.Exit(1)
		}
		process()
	}
}

func process() {
	var file *os.File

	if file_name == "REQUIRED" {
		file = os.Stdin
	} else {
		var err error
		file, err = os.Open(file_name)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer file.Close()

	pokenv := pokenv{registry: registry, checkPath: *flag_check}

	if *flag_machine {
		pokenv.processFile(REG_KEY_MACHINE, file)
	} else {
		pokenv.processFile(REG_KEY_USER, file)
	}
	refreshEnvironment()
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// https://goreleaser.com/cookbooks/using-main.version/
var (
	name    string
	version string
	date    string
	commit  string
)

// flags
type Config struct {
	user    bool
	machine bool
	file    string
	noop    bool
	quiet   bool
	help    bool
	version bool
}

func initFlags() *Config {
	cfg := &Config{}
	flag.BoolVar(&cfg.user, "u", false, "")
	flag.BoolVar(&cfg.user, "user", false, "variables should be set for current user (HKEY_CURRENT_USER) (default)")
	flag.BoolVar(&cfg.machine, "m", false, "")
	flag.BoolVar(&cfg.machine, "machine", false, "variables should be set system wide (HKEY_LOCAL_MACHINE)")
	flag.StringVar(&cfg.file, "f", "stdin", "")
	flag.StringVar(&cfg.file, "file", "stdin", "text file containing the variables to load (default \"stdin\")")
	flag.BoolVar(&cfg.quiet, "q", false, "")
	flag.BoolVar(&cfg.quiet, "quiet", false, "suppress non-error output")
	flag.BoolVar(&cfg.help, "?", false, "")
	flag.BoolVar(&cfg.help, "help", false, "displays this help message")
	flag.BoolVar(&cfg.version, "v", false, "")
	flag.BoolVar(&cfg.version, "version", false, "print version and exit")
	return cfg
}

func main() {
	log.SetFlags(0)
	cfg := initFlags()

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: "+name+` [OPTIONS] [-f file]

Pokenv uses a text file containing the variables to load (default "stdin").
The input format should look like this:

  [GOBIN]
  c:\usr\bin

  [-GOPATH]
  # prefix with dash '-' to remove variable

OPTIONS:
  -u, --user
        set vars for current user (HKEY_CURRENT_USER) (default)
  -m, --machine
        set vars system wide (HKEY_LOCAL_MACHINE)
  -f, --file
        text file containing the variables to load (default "stdin")
  -n, --noop
        do not modify registry
  -q, --quiet
        suppress non-error output
  -?, --help
        display this help message
  -v, --version
        print version and exit

EXAMPLES:`)

		fmt.Fprintln(os.Stderr, "\n  $ "+name+` -f user_variables.ini`)
	}

	flag.Parse()

	if flag.Arg(0) == "version" || cfg.version {
		fmt.Printf("%s %s, built on %s (commit: %s)\n", name, version, date, commit)
		return
	}

	if cfg.help {
		flag.Usage()
		return
	}

	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}

	// TODO: implement noop
	if cfg.noop {
		log.Fatalln("--noop not implemented")
	}
	// TODO: implement quiet
	if cfg.quiet {
		log.Fatalln("--quiet not implemented")
	}

	if !cfg.user && !cfg.machine {
		cfg.user = true // default
	}

	if cfg.user && cfg.machine {
		log.Fatalln("Cannot specify both --user and --machine")
	}

	// process registry
	registry := realRegistry{}
	pokenv := pokenv{registry, cfg}
	var err error
	if cfg.machine {
		err = pokenv.processFile(REG_KEY_MACHINE, cfg.file)
	} else {
		err = pokenv.processFile(REG_KEY_USER, cfg.file)
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

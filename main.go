package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	version       [3]int = [3]int{0, 1, 0}
	optHashVerify string
	optAbs        bool
	optSquash     bool
	optVerbose    bool
	optDebug      bool
)

func printVersion() {
	fmt.Printf("%d.%d.%d\n",
		version[0],
		version[1],
		version[2])
}

func usage(progname string) {
	fmt.Fprintln(os.Stderr, "usage: "+progname+": [<options>] <paths>")
	flag.PrintDefaults()
}

func main() {
	progname := os.Args[0]

	opt_hash_algo := flag.String("hash_algo", SHA256, "Hash algorithm to use")
	opt_hash_verify := flag.String("hash_verify", "", "Message digest to verify in hex string")
	opt_abs := flag.Bool("abs", false, "Use absolute path in output")
	opt_squash := flag.Bool("squash", false, "Enable squashed message digest")
	opt_verbose := flag.Bool("verbose", false, "Enable verbose print")
	opt_debug := flag.Bool("debug", false, "Enable debug print")
	opt_version := flag.Bool("v", false, "Print version and exit")
	opt_help_h := flag.Bool("h", false, "Print usage and exit")

	flag.Parse()
	args := flag.Args()
	optHashVerify = strings.ToLower(*opt_hash_verify)
	optAbs = *opt_abs
	optSquash = *opt_squash
	optVerbose = *opt_verbose
	optDebug = *opt_debug

	if *opt_version {
		printVersion()
		os.Exit(1)
	}

	if *opt_help_h {
		usage(progname)
		os.Exit(1)
	}

	if len(args) < 1 {
		usage(progname)
		os.Exit(1)
	}

	hash_algo := strings.ToLower(*opt_hash_algo)
	if hash_algo == "" {
		fmt.Println("No hash algorithm specified")
		os.Exit(1)
	}
	if optVerbose {
		fmt.Println(hash_algo)
	}

	if newHash(hash_algo) == nil {
		fmt.Println("Unsupported hash algorithm", hash_algo)
		fmt.Println("Available hash algorithm", getAvailableHashAlgo())
		os.Exit(1)
	}

	if optHashVerify != "" {
		var valid bool
		if optHashVerify, valid = isValidHexSum(optHashVerify); !valid {
			fmt.Println("Invalid verify string", optHashVerify)
			os.Exit(1)
		}
	}
	assert(optHashVerify == strings.ToLower(optHashVerify))

	for _, x := range args {
		err := printInput(x, hash_algo)
		if err != nil {
			panic(err)
		}
	}
}

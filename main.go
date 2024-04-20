package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

var (
	version          [3]int = [3]int{0, 4, 4}
	optHashAlgo      string
	optHashVerify    string
	optHashOnly      bool
	optIgnoreDot     bool
	optIgnoreDotDir  bool
	optIgnoreDotFile bool
	optIgnoreSymlink bool
	optFollowSymlink bool
	optAbs           bool
	optSwap          bool
	optSort          bool
	optSquash        bool
	optVerbose       bool
	optDebug         bool
)

func getVersionString() string {
	return fmt.Sprintf("%d.%d.%d", version[0], version[1], version[2])
}

func printVersion() {
	fmt.Println(getVersionString())
}

func usage(progname string) {
	fmt.Fprintln(os.Stderr, "usage: "+progname+": [<options>] <paths>")
	flag.PrintDefaults()
}

func main() {
	progname := path.Base(os.Args[0])

	opt_hash_algo := flag.String("hash_algo", SHA256, "Hash algorithm to use")
	opt_hash_verify := flag.String("hash_verify", "", "Message digest to verify in hex string")
	opt_hash_only := flag.Bool("hash_only", false, "Do not print file paths")
	opt_ignore_dot := flag.Bool("ignore_dot", false, "Ignore entries start with .")
	opt_ignore_dot_dir := flag.Bool("ignore_dot_dir", false, "Ignore directories start with .")
	opt_ignore_dot_file := flag.Bool("ignore_dot_file", false, "Ignore files start with .")
	opt_ignore_symlink := flag.Bool("ignore_symlink", false, "Ignore symbolic links")
	opt_follow_symlink := flag.Bool("follow_symlink", false, "Follow symbolic links unless directory")
	opt_abs := flag.Bool("abs", false, "Print file paths in absolute path")
	opt_swap := flag.Bool("swap", false, "Print file path first in each line")
	opt_sort := flag.Bool("sort", false, "Print sorted file paths")
	opt_squash := flag.Bool("squash", false, "Print squashed message digest instead of per file")
	opt_verbose := flag.Bool("verbose", false, "Enable verbose print")
	opt_debug := flag.Bool("debug", false, "Enable debug print")
	opt_version := flag.Bool("v", false, "Print version and exit")
	opt_help_h := flag.Bool("h", false, "Print usage and exit")

	flag.Parse()
	args := flag.Args()
	optHashAlgo = strings.ToLower(*opt_hash_algo)
	optHashVerify = strings.ToLower(*opt_hash_verify)
	optHashOnly = *opt_hash_only
	optIgnoreDot = *opt_ignore_dot
	optIgnoreDotDir = *opt_ignore_dot_dir
	optIgnoreDotFile = *opt_ignore_dot_file
	optIgnoreSymlink = *opt_ignore_symlink
	optFollowSymlink = *opt_follow_symlink
	optAbs = *opt_abs
	optSwap = *opt_swap
	optSort = *opt_sort
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

	if len(optHashAlgo) == 0 {
		fmt.Println("No hash algorithm specified")
		os.Exit(1)
	}
	if optVerbose {
		fmt.Println(optHashAlgo)
	}

	if newHash(optHashAlgo) == nil {
		fmt.Println("Unsupported hash algorithm", optHashAlgo)
		fmt.Println("Available hash algorithm", getAvailableHashAlgo())
		os.Exit(1)
	}

	if len(optHashVerify) != 0 {
		var valid bool
		if optHashVerify, valid = isValidHexSum(optHashVerify); !valid {
			fmt.Println("Invalid verify string", optHashVerify)
			os.Exit(1)
		}
	}
	assert(optHashVerify == strings.ToLower(optHashVerify))

	if isWindows() {
		fmt.Println("Windows unsupported")
		os.Exit(1)
	}

	if s := getPathSeparator(); s != '/' {
		fmt.Println("Invalid path separator", s)
		os.Exit(1)
	}

	for i, x := range args {
		if err := printInput(x); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if optVerbose && len(args) > 0 && i != len(args)-1 {
			fmt.Println()
		}
	}
}

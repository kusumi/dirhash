package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

var (
	version          [3]int = [3]int{0, 4, 5}
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

	optHashAlgoAddr := flag.String("hash_algo", SHA256, "Hash algorithm to use")
	optHashVerifyAddr := flag.String("hash_verify", "", "Message digest to verify in hex string")
	optHashOnlyAddr := flag.Bool("hash_only", false, "Do not print file paths")
	optIgnoreDotAddr := flag.Bool("ignore_dot", false, "Ignore entries start with .")
	optIgnoreDotDirAddr := flag.Bool("ignore_dot_dir", false, "Ignore directories start with .")
	optIgnoreDotFileAddr := flag.Bool("ignore_dot_file", false, "Ignore files start with .")
	optIgnoreSymlinkAddr := flag.Bool("ignore_symlink", false, "Ignore symbolic links")
	optFollowSymlinkAddr := flag.Bool("follow_symlink", false, "Follow symbolic links unless directory")
	optAbsAddr := flag.Bool("abs", false, "Print file paths in absolute path")
	optSwapAddr := flag.Bool("swap", false, "Print file path first in each line")
	optSortAddr := flag.Bool("sort", false, "Print sorted file paths")
	optSquashAddr := flag.Bool("squash", false, "Print squashed message digest instead of per file")
	optVerboseAddr := flag.Bool("verbose", false, "Enable verbose print")
	optDebugAddr := flag.Bool("debug", false, "Enable debug print")
	optVersionAddr := flag.Bool("v", false, "Print version and exit")
	optHelpAddr := flag.Bool("h", false, "Print usage and exit")

	flag.Parse()
	args := flag.Args()
	optHashAlgo = strings.ToLower(*optHashAlgoAddr)
	optHashVerify = strings.ToLower(*optHashVerifyAddr)
	optHashOnly = *optHashOnlyAddr
	optIgnoreDot = *optIgnoreDotAddr
	optIgnoreDotDir = *optIgnoreDotDirAddr
	optIgnoreDotFile = *optIgnoreDotFileAddr
	optIgnoreSymlink = *optIgnoreSymlinkAddr
	optFollowSymlink = *optFollowSymlinkAddr
	optAbs = *optAbsAddr
	optSwap = *optSwapAddr
	optSort = *optSortAddr
	optSquash = *optSquashAddr
	optVerbose = *optVerboseAddr
	optDebug = *optDebugAddr

	if *optVersionAddr {
		printVersion()
		os.Exit(1)
	}

	if *optHelpAddr {
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

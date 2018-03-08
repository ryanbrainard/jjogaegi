package main

import (
	"flag"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ryanbrainard/jjogaegi/cmd"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/ryanbrainard/jjogaegi/run"

	"golang.org/x/crypto/ssh/terminal"
)

var fVersion = flag.Bool("version", false, "print version information")
var fIn = flag.String("in", "stdin", "filename to read as input")
var fOut = flag.String("out", "stdout", "filename to write to output")
var fParser = flag.String("parser", "prompt", "type of parser for input ["+strings.Join(cmd.Keys(cmd.AppCapabilities.Parsers), "|")+"]")
var fFormatter = flag.String("formatter", "tsv", "type of formatter for output ["+strings.Join(cmd.Keys(cmd.AppCapabilities.Formatters), "|")+"]")
var fHanja = flag.String("hanja", "none", "include hanja [none|parens]")
var fHeader = flag.String("header", "", "header to prepend to output")
var fMediadir = flag.String("mediadir", "", "dir to download media. alternatively set with MEDIA_DIR env.")
var fParallel = flag.Bool("parallel", false, "parallel processing. records may be returned out of order.")
var fLookup = flag.Bool("lookup", false, "look up words in dictionary to enhance item details. always true with prompt parser.")
var fInteractive = flag.Bool("interactive", false, "interactive mode. always true with prompt parser.")

func main() {
	flag.Parse()

	if *fVersion {
		os.Stderr.WriteString("v" + pkg.VERSION + "\n")
		os.Exit(0)
	}

	var in io.ReadCloser
	var err error
	switch *fIn {
	case "stdin":
		in = os.Stdin
		if *fInteractive && *fParser != "prompt" {
			os.Stderr.WriteString("Interactive mode cannot be used with " + *fParser + " parser on stdin. Set -in option or do not set parser. Run `jjogaegi -help` for details.\n")
			os.Exit(4)
		}
	default:
		if in, err = os.Open(*fIn); err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(2)
		}
		defer in.Close()
	}

	var out io.WriteCloser
	switch *fOut {
	case "stdout":
		if (*fParser == "prompt" || *fInteractive) && terminal.IsTerminal(int(os.Stdout.Fd())) {
			os.Stderr.WriteString("Set -out option or redirect outout when using interactive mode.  Run `jjogaegi -help` for details.\n")
			os.Exit(10)
		}
		out = os.Stdout
	default:
		if out, err = os.Create(*fOut); err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(3)
		}
		defer out.Close()
	}

	err = run.Run(
		in,
		out,
		cmd.ParseOptParser(*fParser),
		cmd.ParseOptFormatter(*fFormatter),
		map[string]string{
			pkg.OPT_HANJA:       cmd.ParseOptHanja(*fHanja),
			pkg.OPT_HEADER:      *fHeader,
			pkg.OPT_MEDIADIR:    *fMediadir,
			pkg.OPT_PARALLEL:    strconv.FormatBool(*fParallel),
			pkg.OPT_LOOKUP:      strconv.FormatBool(*fLookup),
			pkg.OPT_INTERACTIVE: strconv.FormatBool(*fInteractive),
		},
	)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

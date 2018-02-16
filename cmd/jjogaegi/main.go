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
)

var fIn = flag.String("in", "stdin", "filename to read as input")
var fOut = flag.String("out", "stdout", "filename to write to as output")
var fParser = flag.String("parser", "list", "type of parser for input ["+strings.Join(cmd.Keys(cmd.AppCapabilities.Parsers), "|")+"]")
var fFormatter = flag.String("formatter", "tsv", "type of formatter for output ["+strings.Join(cmd.Keys(cmd.AppCapabilities.Formatters), "|")+"]")
var fHanja = flag.String("hanja", "none", "include hanja [none|parens]")
var fHeader = flag.String("header", "", "header to prepend to output")
var fMediadir = flag.String("mediadir", "", "dir to download media")
var fParallel = flag.Bool("parallel", false, "process in parallel (records may be returned out of order)")
var fLookup = flag.Bool("lookup", false, "look up words in dictionary") // TODO: expand to allow custom lookup
var fInteractive = flag.Bool("interactive", false, "allow interactive prompting")

func main() {
	flag.Parse()

	var in io.ReadCloser
	var err error
	switch *fIn {
	case "stdin":
		if *fInteractive {
			os.Stderr.WriteString("Cannot use interactive mode while reading from stdin. Set -in option for input from a file.\n")
			os.Exit(4)
		}
		in = os.Stdin
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

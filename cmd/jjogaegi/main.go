package main

import (
	"flag"
	"fmt"
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
var fHeader = flag.String("header", "", "header to prepend to output")
var fParallel = flag.Bool("parallel", false, "parallel processing. records may be returned out of order.")
var fLookup = flag.Bool("lookup", false, "look up words in dictionary to enhance item details. always true with prompt parser.")
var fInteractive = flag.Bool("interactive", false, "interactive mode. always true with prompt parser.")

func main() {
	flag.Usage = func() {
		scriptName := os.Args[0]

		fmt.Fprintf(os.Stderr, "%s - Korean vocabulary parser-formatter\n\n", scriptName)
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", scriptName)
		fmt.Fprintf(os.Stderr, "       Parses input interactively, from stdin, or from a file specificed with the -in option.\n")
		fmt.Fprintf(os.Stderr, "       Formats output to stdout or to a file specificed with the -out option.\n")
		fmt.Fprintf(os.Stderr, "       Set environment for global configuration and options for per run configuration.\n\n")
		fmt.Fprintf(os.Stderr, "Environment:\n")
		fmt.Fprintf(os.Stderr, "  KRDICT_API_KEY: Dictionary API key to enable word lookups.\n")
		fmt.Fprintf(os.Stderr, "                  For registration, see https://krdict.korean.go.kr/openApi/openApiInfo\n")
		fmt.Fprintf(os.Stderr, "  MEDIA_DIR:      Directory to download images and audio.\n")
		fmt.Fprintf(os.Stderr, "                  For use with Anki, see https://apps.ankiweb.net/docs/manual.html#files\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%s/%s\thttps://github.com/ryanbrainard/jjogaegi\n", scriptName, pkg.VERSION)

	}

	flag.Parse()

	if *fVersion {
		fmt.Printf("v%s\n", pkg.VERSION)
		os.Exit(0)
	}

	var in io.ReadCloser
	var err error
	switch *fIn {
	case "stdin":
		in = os.Stdin
		if *fInteractive && *fParser != "prompt" {
			fmt.Fprintf(os.Stderr, "Interactive mode cannot be used with %s parser on stdin. Set -in option or do not set parser. Run `jjogaegi -help` for details.\n", *fParser)
			os.Exit(4)
		}
	default:
		if in, err = os.Open(*fIn); err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
			os.Exit(2)
		}
		defer in.Close()
	}

	var out io.WriteCloser
	switch *fOut {
	case "stdout":
		if (*fParser == "prompt" || *fInteractive) && terminal.IsTerminal(int(os.Stdout.Fd())) {
			fmt.Fprintf(os.Stderr, "Set -out option or redirect outout when using interactive mode.  Run `jjogaegi -help` for details.\n")
			os.Exit(10)
		}
		out = os.Stdout
	default:
		if out, err = os.Create(*fOut); err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
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
			pkg.OPT_HEADER:      *fHeader,
			pkg.OPT_PARALLEL:    strconv.FormatBool(*fParallel),
			pkg.OPT_LOOKUP:      strconv.FormatBool(*fLookup),
			pkg.OPT_INTERACTIVE: strconv.FormatBool(*fInteractive),
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}
}

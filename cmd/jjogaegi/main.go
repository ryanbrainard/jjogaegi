package main

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/ryanbrainard/jjogaegi/cmd"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/ryanbrainard/jjogaegi/run"
)

var fParser = flag.String("parser", "list", "type of parser for input ["+strings.Join(cmd.Keys(cmd.AppCapabilities.Parsers), "|")+"]")
var fFormatter = flag.String("formatter", "tsv", "type of formatter for output ["+strings.Join(cmd.Keys(cmd.AppCapabilities.Formatters), "|")+"]")
var fHanja = flag.String("hanja", "none", "include hanja [none|parens]")
var fHeader = flag.String("header", "", "header to prepend to output")
var fMediadir = flag.String("mediadir", "", "dir to download media")
var fParallel = flag.Bool("parallel", false, "process in parallel (records may be returned out of order)")
var fKrDictLookup = flag.Bool("lookup", false, "look up words in dictionary") // TODO: expand to allow custom lookup

func main() {
	flag.Parse()
	err := run.Run(
		os.Stdin,
		os.Stdout,
		cmd.ParseOptParser(*fParser),
		cmd.ParseOptFormatter(*fFormatter),
		map[string]string{
			pkg.OPT_HANJA:         cmd.ParseOptHanja(*fHanja),
			pkg.OPT_HEADER:        *fHeader,
			pkg.OPT_MEDIADIR:      *fMediadir,
			pkg.OPT_PARALLEL:      strconv.FormatBool(*fParallel),
			pkg.OPT_KRDICT_LOOKUP: strconv.FormatBool(*fKrDictLookup),
		},
	)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

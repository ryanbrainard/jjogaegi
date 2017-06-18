package main

import (
	"flag"
	"os"
	"strings"

	"github.com/ryanbrainard/jjogaegi/cmd"
	"github.com/ryanbrainard/jjogaegi/pkg"
)

var fParser = flag.String("parser", "list", "type of parser for input [" + strings.Join(cmd.Keys(cmd.AppCapabilities.Parsers), "|") + "]")
var fFormatter = flag.String("formatter", "tsv", "type of formatter for output [" + strings.Join(cmd.Keys(cmd.AppCapabilities.Formatters), "|") + "]")
var fHanja = flag.String("hanja", "none", "include hanja [none|parens]")
var fHeader = flag.String("header", "", "header to prepend to output")
var fAudiodir = flag.String("audiodir", "", "dir to download audio")

func main() {
	flag.Parse()
	err := pkg.Run(
		os.Stdin,
		os.Stdout,
		cmd.ParseOptParser(*fParser),
		cmd.ParseOptFormatter(*fFormatter),
		map[string]string{
			pkg.OPT_HANJA:    cmd.ParseOptHanja(*fHanja),
			pkg.OPT_HEADER:   *fHeader,
			pkg.OPT_AUDIODIR: *fAudiodir,
		},
	)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}

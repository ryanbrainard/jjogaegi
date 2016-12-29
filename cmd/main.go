package main

import (
	"os"
	"flag"
	"github.com/ryanbrainard/jjogaegi"
	"github.com/ryanbrainard/jjogaegi/parsers"
	"github.com/ryanbrainard/jjogaegi/formatters"
)

var flagParser = flag.String("parser", "list", "type of parser for input [list|table]")
var flagFormatter = flag.String("formatter", "tsv", "type of formatter for output [tsv|csv]")

func main() {
	flag.Parse()
	jjogaegi.Run(
		os.Stdin,
		parser(),
		formatter(),
		os.Stdout,
	)
}

func parser() jjogaegi.ParseFunc {
	switch *flagParser {
	case "list":
		return parsers.ParseList
	case "table":
		return parsers.ParseTable
	default:
		panic("Unknown parser")
	}
}

func formatter() jjogaegi.FormatFunc {
	switch *flagFormatter {
	case "tsv":
		return formatters.FormatTSV
	case "csv":
		return formatters.FormatCSV
	default:
		panic("Unknown formatter")
	}
}

package main

import (
	"flag"
	"github.com/ryanbrainard/jjogaegi/formatters"
	"github.com/ryanbrainard/jjogaegi/parsers"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"os"
)

var fParser = flag.String("parser", "list", "type of parser for input [list|naver-table|naver-json]")
var fFormatter = flag.String("formatter", "tsv", "type of formatter for output [tsv|csv]")
var fHanja = flag.String("hanja", "none", "include hanja [none|parens|sep]")
var fHeader = flag.String("header", "", "header to prepend to output")

func main() {
	flag.Parse()
	pkg.Run(
		os.Stdin,
		os.Stdout,
		parser(),
		formatter(),
		options(),
	)
}

func parser() pkg.ParseFunc {
	switch *fParser {
	case "list":
		return parsers.ParseList
	case "naver-table":
		return parsers.ParseNaverTable
	case "naver-json":
		return parsers.ParseNaverJSON
	default:
		panic("Unknown parser")
	}
}

func formatter() pkg.FormatFunc {
	switch *fFormatter {
	case "tsv":
		return formatters.FormatTSV
	case "csv":
		return formatters.FormatCSV
	default:
		panic("Unknown formatter")
	}
}

func options() map[string]string {
	return map[string]string{
		pkg.OPT_HANJA: hanja(),
		pkg.OPT_HEADER: *fHeader,
	}
}

func hanja() string {
	switch *fHanja {
	case pkg.OPT_HANJA_NONE, pkg.OPT_HANJA_PARENTHESIS, pkg.OPT_HANJA_SEPARATE:
		return *fHanja
	default:
		panic("Unknown hanja option")
	}
}

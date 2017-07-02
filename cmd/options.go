package cmd

import (
	"github.com/ryanbrainard/jjogaegi/formatters"
	"github.com/ryanbrainard/jjogaegi/parsers"
	"github.com/ryanbrainard/jjogaegi/pkg"
)

type Capabilities struct {
	Parsers    map[string]string
	Formatters map[string]string
}

var AppCapabilities = Capabilities {
	Parsers: map[string]string{
		"list":         "List",
		"naver-table":  "Naver Table",
		"naver-json":   "Naver JSON",
		"krdict-xml":   "KR Dict XML",
		"memrise-list": "Memrise List",
	},
	Formatters: map[string]string{
		"tsv": "TSV: Tab-Separated Values",
		"csv": "CSV: Comma-Separated Values",
	},
}

func Keys(m map[string]string) []string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func ParseOptParser(s string) pkg.ParseFunc {
	switch s {
	case "list":
		return parsers.ParseList
	case "naver-table":
		return parsers.ParseNaverTable
	case "naver-json":
		return parsers.ParseNaverJSON
	case "krdict-xml":
		return parsers.ParseKrDictXML
	case "memrise-list":
		return parsers.ParseMemriseList
	default:
		return nil
	}
}

func ParseOptFormatter(s string) pkg.FormatFunc {
	switch s {
	case "tsv":
		return formatters.FormatTSV
	case "csv":
		return formatters.FormatCSV
	default:
		return nil
	}
}

func ParseOptHanja(s string) string {
	switch s {
	case pkg.OPT_HANJA_NONE, pkg.OPT_HANJA_PARENTHESIS:
		return s
	default:
		return ""
	}
}

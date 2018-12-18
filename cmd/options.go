package cmd

import (
	"os"

	"github.com/ryanbrainard/jjogaegi/formatters"
	"github.com/ryanbrainard/jjogaegi/parsers"
	"github.com/ryanbrainard/jjogaegi/pkg"
)

type Capabilities struct {
	Parsers    map[string]string
	Formatters map[string]string
}

var AppCapabilities = Capabilities{
	Parsers: map[string]string{
		"prompt":         "Interactive Prompt",
		"tsv":            "TSV: Tab-Separated Values",
		"list":           "Hangul-English Space-Separated List",
		"naver-wordbook": "Naver Wordbook Print Out HTML",
		"krdict-xml":     "KR Dict XML",
		"memrise-list":   "Memrise List",
	},
	Formatters: map[string]string{
		"json": "JSON",
		"tsv":  "TSV: Tab-Separated Values",
		"csv":  "CSV: Comma-Separated Values",
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
	case "prompt":
		return parsers.NewInteractivePrompt(os.Stderr)
	case "tsv":
		return parsers.ParseTSV
	case "list":
		return parsers.ParseList
	case "naver-wordbook":
		return parsers.ParseNaverWordbook
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
	case "json":
		return formatters.FormatJSON
	case "tsv":
		return formatters.FormatTSV
	case "csv":
		return formatters.FormatCSV
	default:
		return nil
	}
}

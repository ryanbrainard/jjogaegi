package cmd

import (
	"os"

	"go.ryanbrainard.com/jjogaegi/formatters"
	"go.ryanbrainard.com/jjogaegi/parsers"
	"go.ryanbrainard.com/jjogaegi/pkg"
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
		"naver-wordbook-html": "Naver Wordbook Print Out HTML",
		"naver-wordbook-json": "Naver Wordbook JSON", // example URL: https://learn.dict.naver.com/gateway-api/enkodict/mywordbook/word/list/search?page_size=20&wbId=a41dd7534bd04c02b50942f6fa935f51&qt=0&st=0&hasBookmark=true&bookmarkWordId=f855470720c6a72b7108df104cd461fa&domain=naver
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
	case "naver-wordbook-html":
		return parsers.ParseNaverWordbookHTML
	case "naver-wordbook-json":
		return parsers.ParseNaverWordbookJSON
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

package cmd

import (
	"github.com/ryanbrainard/jjogaegi/formatters"
	"github.com/ryanbrainard/jjogaegi/parsers"
	"github.com/ryanbrainard/jjogaegi/pkg"
)

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

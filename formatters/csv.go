package formatters

import (
	"encoding/csv"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"io"
)

func FormatCSV(items <-chan *pkg.Item, w io.Writer) {
	formatXSV(items, w, ',')
}

func formatXSV(items <-chan *pkg.Item, w io.Writer, delim rune) {
	cw := csv.NewWriter(w)
	cw.Comma = delim
	for item := range items {
		cw.Write([]string{
			mergeTermSubTerm(item),
			item.Def,
		})
	}
	cw.Flush()
}

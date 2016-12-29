package formatters

import (
	"encoding/csv"
	"io"
	"github.com/ryanbrainard/jjogaegi"
)

func FormatCSV(items <-chan *jjogaegi.Item, w io.Writer) {
	formatXSV(items, w, ',')
}

func formatXSV(items <-chan *jjogaegi.Item, w io.Writer, delim rune) {
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

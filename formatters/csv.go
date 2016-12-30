package formatters

import (
	"encoding/csv"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"io"
)

func FormatCSV(items <-chan *pkg.Item, w io.Writer, options map[string]string) {
	formatXSV(items, w, options, ',')
}

func formatXSV(items <-chan *pkg.Item, w io.Writer, options map[string]string, delim rune) {
	cw := csv.NewWriter(w)
	cw.Comma = delim
	for item := range items {
		cw.Write(append(formatHangulHanja(item, options), item.Def))
	}
	cw.Flush()
}

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
	writeHeader(w, options)
	cw := csv.NewWriter(w)
	cw.Comma = delim
	for item := range items {
		var firstExample pkg.Example
		if len(item.Examples) > 0 {
			firstExample = item.Examples[0]
		}

		cw.Write(append(formatHangulHanja(item, options), item.Def, firstExample.Korean, firstExample.English))
	}
	cw.Flush()
}

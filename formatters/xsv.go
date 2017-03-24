package formatters

import (
	"encoding/csv"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"io"
)

func formatXSV(items <-chan *pkg.Item, w io.Writer, options map[string]string, delim rune) {
	writeHeader(w, options)
	cw := csv.NewWriter(w)
	cw.Comma = delim
	for item := range items {
		var firstExample pkg.Translation
		if len(item.Examples) > 0 {
			firstExample = item.Examples[0]
		}

		var secondExample pkg.Translation
		if len(item.Examples) > 1 {
			secondExample = item.Examples[1]
		}

		var audioTag string
		if item.AudioURL != "" {
			audioTag = "[sound:" + item.AudioURL + "]"
		}

		cw.Write([]string{
			item.Id,
			formatHangulHanja(item, options),
			item.Hanja,
			item.Pronunciation,
			audioTag,
			item.Def.Korean,
			item.Def.English,
			item.Antonym,
			firstExample.Korean,
			firstExample.English,
			secondExample.Korean,
			secondExample.English,
		})
	}
	cw.Flush()
}

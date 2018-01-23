package formatters

import (
	"encoding/csv"
	"io"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

func formatXSV(items <-chan *pkg.Item, w io.Writer, options map[string]string, delim rune) error {
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

		audioTag, err := formatAudioTag(item, options)
		if err != nil {
			return err
		}

		imageTag, err := formatImageTag(item, options)
		if err != nil {
			return err
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
			imageTag,
		})
	}
	cw.Flush()
	return nil
}

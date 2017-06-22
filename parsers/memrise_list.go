package parsers

import (
	"io"

	"bufio"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

func ParseMemriseList(r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	scanner := bufio.NewScanner(r)

	hangul := ""
	for scanner.Scan() {
		// hangul is on first line; def is on second line,
		// but we scan one line at a time.
		line := scanner.Text()
		if hangul == "" {
			hangul = line
			continue
		}
		def := line

		items <- &pkg.Item{
			Id:     sanitize(hangul),
			Hangul: sanitize(hangul),
			Def: pkg.Translation{
				English: sanitize(def),
			},
		}

		hangul = ""
	}

	return nil
}

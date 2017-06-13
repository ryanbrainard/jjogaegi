package parsers

import (
	"bufio"
	"io"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

func ParseList(r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	scanner := bufio.NewScanner(r)

	line := ""
	for scanner.Scan() {
		line = line + scanner.Text()
		hangul, def := splitHangul(line)
		if !isMultiline(options) || (len(sanitize(hangul)) > 0 && len(sanitize(def)) > 0) {
			items <- &pkg.Item{
				Id: sanitize(hangul),
				Hangul: sanitize(hangul),
				Def: pkg.Translation{
					English: sanitize(def),
				},
			}
			line = ""
		}
	}

	return nil
}

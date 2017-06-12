package parsers

import (
	"bufio"
	"io"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

func ParseList(r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		hangul, def := splitHangul(line)
		items <- &pkg.Item{
			Id: sanitize(hangul),
			Hangul: sanitize(hangul),
			Def: pkg.Translation{
				English: sanitize(def),
			},
		}
	}

	return nil
}

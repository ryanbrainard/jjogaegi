package parsers

import (
	"bufio"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"io"
)

func ParseList(r io.Reader, items chan<- *pkg.Item, options map[string]string) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		hangul, def := splitHangul(line)
		items <- &pkg.Item{
			Hangul: sanitize(hangul),
			Def: pkg.Translation{
				English: sanitize(def),
			},
		}
	}

	close(items)
}

package parsers

import (
	"context"
	"io"

	"bufio"

	"ryanbrainard.com/jjogaegi/pkg"
)

func ParseMemriseList(ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	scanner := bufio.NewScanner(r)

	hangul := ""
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// hangul is on first line; def is on second line,
		// but we scan one line at a time.
		line := scanner.Text()
		if hangul == "" {
			hangul = line
			continue
		}
		def := line

		items <- &pkg.Item{
			Hangul: sanitize(hangul),
			Def: pkg.Translation{
				English: sanitize(def),
			},
		}

		hangul = ""
	}

	return nil
}

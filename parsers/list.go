package parsers

import (
	"bufio"
	"context"
	"io"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

func ParseList(ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()
		hangul, def := splitHangul(line)
		items <- &pkg.Item{
			Hangul: sanitize(hangul),
			Def: pkg.Translation{
				English: sanitize(def),
			},
		}
	}

	return nil
}

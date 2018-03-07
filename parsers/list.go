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

		items <- parseLineItem(scanner.Text())
	}

	return nil
}

func parseLineItem(line string) *pkg.Item {
	hangul, def := splitHangul(line)
	return &pkg.Item{
		Hangul: sanitize(hangul),
		Def: pkg.Translation{
			English: sanitize(def),
		},
	}
}

package parsers

import (
	"bufio"
	"context"
	"io"

	"go.ryanbrainard.com/jjogaegi/pkg"
)

func ParseList(ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		item := parseLineItem(scanner.Text())
		if item.Hangul == "" {
			continue
		}
		items <- item
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

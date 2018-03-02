package parsers

import (
	"bufio"
	"context"
	"io"
	"strings"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

func ParseNaverTable(ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	i := 0
	rawTerms := []string{}
	scanner := bufio.NewScanner(r)
	scanner.Split(SplitDefs)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()

		if line == "" {
			continue
		}

		if len(rawTerms) == 0 {
			rawTerms = strings.Split(line, " ")
			continue
		}

		hangulTerm, hanjaTerm := splitHangul(rawTerms[i])

		items <- &pkg.Item{
			Hangul: sanitize(string(hangulTerm)),
			Hanja:  sanitize(string(hanjaTerm)),
			Def: pkg.Translation{
				English: sanitize(string(line)),
			},
		}

		i++
	}

	return nil
}

func SplitDefs(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF {
		if len(data) == 0 {
			return 0, nil, nil
		}
		return len(data), data, nil
	}

	oneByte := []byte("1")[0]
	for i, d := range data {
		if d == oneByte {
			return i + 1, data[0:i], nil
		}
	}

	return 0, nil, nil
}

package parsers

import (
	"bufio"
	"github.com/ryanbrainard/jjogaegi"
	"io"
	"strings"
)

func ParseTable(r io.Reader, items chan<- *jjogaegi.Item) {
	i := 0
	rawTerms := []string{}
	scanner := bufio.NewScanner(r)
	scanner.Split(SplitDefs)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		if len(rawTerms) == 0 {
			rawTerms = strings.Split(line, " ")
			continue
		}

		rawTerm := rawTerms[i]
		hangeulTerm := []rune{}
		hanjaTerm := []rune{}
		for j, t := range rawTerm {
			if hasHangeul(rawTerm[j:]) {
				hangeulTerm = append(hangeulTerm, t)
			} else {
				hanjaTerm = append(hanjaTerm, t)
			}
		}

		items <- &jjogaegi.Item{
			Term:    sanitize(string(hangeulTerm)),
			SubTerm: sanitize(string(hanjaTerm)),
			Def:     sanitize(string(line)),
		}

		i++
	}

	close(items)
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

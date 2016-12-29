package parsers

import (
	"bufio"
	"github.com/ryanbrainard/jjogaegi"
	"io"
)

func ParseList(r io.Reader, items chan<- *jjogaegi.Item) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		term, def := splitHangeul(line)
		items <- &jjogaegi.Item{
			Term: sanitize(term),
			Def:  sanitize(def),
		}
	}

	close(items)
}

package parsers

import (
	"bufio"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"io"
)

func ParseList(r io.Reader, items chan<- *pkg.Item) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		term, def := splitHangeul(line)
		items <- &pkg.Item{
			Term: sanitize(term),
			Def:  sanitize(def),
		}
	}

	close(items)
}

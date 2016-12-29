package parsers

import (
	"bufio"
	"io"
	"github.com/ryanbrainard/jjogaegi"
)

func ParseList(r io.Reader, items chan<- *jjogaegi.Item) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		term := []rune{}
		def := []rune{}

		for i, c := range line {
			if isHeader(term, c) {
				continue
			} else if hasHangeul(line[i:]) {
				term = append(term, c)
			} else {
				def = append(def, c)
			}
		}

		items <- &jjogaegi.Item{
			Term: sanitize(string(term)),
			Def: sanitize(string(def)),
		}
	}

	close(items)
}

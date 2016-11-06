package jjogaegi

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

var hangeulRange = []rune("가힣")
var cutSet = " “”."

func Parse(r io.Reader, w io.Writer) {
	in := bufio.NewScanner(r)
	for in.Scan() {
		line := in.Text()
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

		sTerm := sanitize(term)
		sDef := sanitize(def)
		if len(sTerm) > 0 && len(sDef) > 0 {
			w.Write([]byte(fmt.Sprintf("%s\t%s\n", sTerm, sDef)))
		}
	}
}

func isHangeul(r rune) bool {
	return r >= hangeulRange[0] && r <= hangeulRange[1]
}

func hasHangeul(s string) bool {
	for _, r := range s {
		if isHangeul(r) {
			return true
		}
	}
	return false
}

func isHeader(term []rune, r rune) bool {
	return len(term) == 0 && !isHangeul(r)
}

func sanitize(rs []rune) string {
	return strings.Trim(string(rs), cutSet)
}

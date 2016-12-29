package parsers

import (
	"strings"
)

var hangeulRange = []rune("가힣")
var cutSet = " :“”.\n"

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

func splitHangeul(s string) (hangeul string, rest string) {
	hangeulRunes := []rune{}
	restRunes := []rune{}

	for i, c := range s {
		if isHeader(hangeulRunes, c) {
			continue
		} else if hasHangeul(s[i:]) {
			hangeulRunes = append(hangeulRunes, c)
		} else {
			restRunes = append(restRunes, c)
		}
	}

	return string(hangeulRunes), string(restRunes)
}

func isHeader(term []rune, r rune) bool {
	return len(term) == 0 && !isHangeul(r)
}

func sanitize(s string) string {
	return strings.TrimSpace(strings.Trim(s, cutSet))
}

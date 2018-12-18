package parsers

import (
	"strings"
)

var hangulRange = []rune("가힣")
var cutSet = " :“”.\n"

func isHangul(r rune) bool {
	return r >= hangulRange[0] && r <= hangulRange[1]
}

func hasHangul(s string) bool {
	for _, r := range s {
		if isHangul(r) {
			return true
		}
	}
	return false
}

func splitHangul(s string) (hangul string, rest string) {
	hangulRunes := []rune{}
	restRunes := []rune{}

	for i, c := range s {
		if isHeader(hangulRunes, c) {
			continue
		} else if hasHangul(s[i:]) {
			hangulRunes = append(hangulRunes, c)
		} else {
			restRunes = append(restRunes, c)
		}
	}

	return string(hangulRunes), string(restRunes)
}

func splitHangulReverse(s string) (rest string, hangul string) {
	restRunes := []rune{}
	hangulRunes := []rune{}

	for _, c := range s {
		if isHangul(c) || hasHangul(string(hangulRunes)) {
			hangulRunes = append(hangulRunes, c)
		} else {
			restRunes = append(restRunes, c)
		}
	}

	return string(restRunes), string(hangulRunes)
}

func isHeader(term []rune, r rune) bool {
	return len(term) == 0 && !isHangul(r)
}

func sanitize(s string) string {
	return strings.TrimSpace(strings.Trim(s, cutSet))
}

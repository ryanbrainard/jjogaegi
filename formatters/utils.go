package formatters

import (
	"github.com/ryanbrainard/jjogaegi/pkg"
)

func mergeTermSubTerm(item *pkg.Item) string {
	s := item.Hangul
	if len(item.Hanja) > 0 {
		s += " (" + item.Hanja + ")"
	}
	return s
}

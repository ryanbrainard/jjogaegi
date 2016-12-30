package formatters

import (
	"github.com/ryanbrainard/jjogaegi/pkg"
)

func mergeTermSubTerm(item *pkg.Item) string {
	s := item.Term
	if len(item.SubTerm) > 0 {
		s += " (" + item.SubTerm + ")"
	}
	return s
}

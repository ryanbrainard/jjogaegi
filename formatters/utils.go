package formatters

import (
	"github.com/ryanbrainard/jjogaegi"
)

func mergeTermSubTerm(item *jjogaegi.Item) string {
	s := item.Term
	if len(item.SubTerm) > 0 {
		s += " (" + item.SubTerm + ")"
	}
	return s
}

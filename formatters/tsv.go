package formatters

import (
	"io"
	"fmt"
	"github.com/ryanbrainard/jjogaegi"
	"strings"
)

func FormatTSV(items <-chan *jjogaegi.Item, w io.Writer) {
	for item := range items {
		w.Write([]byte(fmt.Sprintf("%s\t%s\n",
			sanitizeTabs(mergeTermSubTerm(item)),
			sanitizeTabs(item.Def))))
	}
}

func sanitizeTabs(s string) string {
	return strings.Replace(s, "\t", "", -1)
}

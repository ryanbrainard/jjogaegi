package formatters

import (
	"io"
	"fmt"
	"github.com/ryanbrainard/jjogaegi"
	"strings"
)

func FormatCSV(items <-chan *jjogaegi.Item, w io.Writer) {
	for item := range items {
		w.Write([]byte(fmt.Sprintf("\"%s\",\"%s\"\n",
			sanitizeSigleQuotes(mergeTermSubTerm(item)),
			sanitizeSigleQuotes(item.Def))))
	}
}


func sanitizeSigleQuotes(s string) string {
	return strings.Replace(s, "'", "\"", -1)
}


package formatters

import (
	"io"
	"fmt"
	"github.com/ryanbrainard/jjogaegi"
)

func FormatCSV(items <-chan *jjogaegi.Item, w io.Writer) {
	for item := range items {
		w.Write([]byte(fmt.Sprintf("%s,%s\n", item.Term, item.Def))) // TODO: subterm
	}
}

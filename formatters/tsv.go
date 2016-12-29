package formatters

import (
	"io"
	"fmt"
	"github.com/ryanbrainard/jjogaegi"
)

func FormatTSV(items <-chan *jjogaegi.Item, w io.Writer) {
	for item := range items {
		w.Write([]byte(fmt.Sprintf("%s\t%s\n", item.Term, item.Def))) // TODO: subterm
	}
}

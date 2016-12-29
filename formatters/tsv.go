package formatters

import (
	"io"
	"github.com/ryanbrainard/jjogaegi"
)

func FormatTSV(items <-chan *jjogaegi.Item, w io.Writer) {
	formatXSV(items, w, '\t')
}

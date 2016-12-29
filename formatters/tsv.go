package formatters

import (
	"github.com/ryanbrainard/jjogaegi"
	"io"
)

func FormatTSV(items <-chan *jjogaegi.Item, w io.Writer) {
	formatXSV(items, w, '\t')
}

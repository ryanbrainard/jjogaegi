package formatters

import (
	"github.com/ryanbrainard/jjogaegi/pkg"
	"io"
)

func FormatTSV(items <-chan *pkg.Item, w io.Writer) {
	formatXSV(items, w, '\t')
}

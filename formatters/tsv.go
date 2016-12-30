package formatters

import (
	"github.com/ryanbrainard/jjogaegi/pkg"
	"io"
)

func FormatTSV(items <-chan *pkg.Item, w io.Writer, options map[string]string) {
	formatXSV(items, w, options, '\t')
}

package formatters

import (
	"github.com/ryanbrainard/jjogaegi/pkg"
	"io"
)

func FormatCSV(items <-chan *pkg.Item, w io.Writer, options map[string]string) {
	formatXSV(items, w, options, ',')
}

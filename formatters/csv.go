package formatters

import (
	"io"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

func FormatCSV(items <-chan *pkg.Item, w io.Writer, options map[string]string) error {
	return formatXSV(items, w, options, ',')
}

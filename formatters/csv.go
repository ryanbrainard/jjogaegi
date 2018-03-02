package formatters

import (
	"context"
	"io"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

func FormatCSV(ctx context.Context, items <-chan *pkg.Item, w io.Writer, options map[string]string) error {
	return formatXSV(ctx, items, w, options, ',')
}

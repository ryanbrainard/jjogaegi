package formatters

import (
	"context"
	"io"

	"go.ryanbrainard.com/jjogaegi/pkg"
)

func FormatTSV(ctx context.Context, items <-chan *pkg.Item, w io.Writer, options map[string]string) error {
	return formatXSV(ctx, items, w, options, '\t')
}

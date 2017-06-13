package parsers

import (
	"io"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

func ParseMemriseList(r io.Reader, items chan <- *pkg.Item, options map[string]string) error {
	options[pkg.OPT_LINEBREAK] = pkg.OPT_LINEBREAK_MULTI
	return ParseList(r, items, options)
}

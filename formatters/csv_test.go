package formatters

import (
	"bytes"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFormatCSV(t *testing.T) {
	items := make(chan *pkg.Item, 1)
	item := &pkg.Item{Hangul: "처리", Hanja: "處理", Def: "handling"}
	out := new(bytes.Buffer)
	items <- item
	close(items)

	FormatCSV(items, out, map[string]string{})

	assert.Equal(t, "처리,handling\n", out.String())
}

func TestFormatCSV_HanjaMerge(t *testing.T) {
	items, out := setupTestFormat()
	FormatCSV(items, out, map[string]string{pkg.OPT_HANJA: pkg.OPT_HANJA_PARENTHESIS})
	assert.Equal(t, "처리 (處理),handling\n", out.String())
}

func TestFormatCSV_HanjaSeparate(t *testing.T) {
	items, out := setupTestFormat()
	FormatCSV(items, out, map[string]string{pkg.OPT_HANJA: pkg.OPT_HANJA_SEPARATE})
	assert.Equal(t, "처리,處理,handling\n", out.String())
}

func setupTestFormat() (<-chan *pkg.Item, *bytes.Buffer) {
	items := make(chan *pkg.Item, 1)
	item := &pkg.Item{Hangul: "처리", Hanja: "處理", Def: "handling"}
	items <- item
	close(items)
	return items, new(bytes.Buffer)
}

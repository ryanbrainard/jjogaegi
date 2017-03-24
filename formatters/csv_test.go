package formatters

import (
	"bytes"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatCSV(t *testing.T) {
	items, out := setupTestFormat()
	FormatCSV(items, out, map[string]string{})
	assert.Equal(t, ",처리,處理,,,,handling,,k,e,,\n", out.String())
}

func TestFormatCSV_Header(t *testing.T) {
	header := "tag: jjogaegi"
	items, out := setupTestFormat()
	FormatCSV(items, out, map[string]string{pkg.OPT_HEADER: header})
	assert.Equal(t, header+"\n,처리,處理,,,,handling,,k,e,,\n", out.String())
}

func TestFormatCSV_HanjaMerge(t *testing.T) {
	items, out := setupTestFormat()
	FormatCSV(items, out, map[string]string{pkg.OPT_HANJA: pkg.OPT_HANJA_PARENTHESIS})
	assert.Equal(t, ",처리 (處理),處理,,,,handling,,k,e,,\n", out.String())
}

func setupTestFormat() (<-chan *pkg.Item, *bytes.Buffer) {
	items := make(chan *pkg.Item, 1)
	item := &pkg.Item{Hangul: "처리", Hanja: "處理", Def: pkg.Translation{English: "handling"}, Examples: []pkg.Translation{{English: "e", Korean: "k"}}}
	items <- item
	close(items)
	return items, new(bytes.Buffer)
}

package formatters

import (
	"bytes"
	"testing"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

func TestFormatCSV(t *testing.T) {
	items, out := setupTestFormat()
	err := FormatCSV(items, out, map[string]string{})
	assert.Nil(t, err)
	assert.Equal(t, ",처리,處理,,,,handling,,k,e,,,\n", out.String())
}

func TestFormatCSV_Header(t *testing.T) {
	header := "tag: jjogaegi"
	items, out := setupTestFormat()
	err := FormatCSV(items, out, map[string]string{pkg.OPT_HEADER: header})
	assert.Nil(t, err)
	assert.Equal(t, header+"\n,처리,處理,,,,handling,,k,e,,,\n", out.String())
}

func TestFormatCSV_HanjaMerge(t *testing.T) {
	items, out := setupTestFormat()
	err := FormatCSV(items, out, map[string]string{pkg.OPT_HANJA: pkg.OPT_HANJA_PARENTHESIS})
	assert.Nil(t, err)
	assert.Equal(t, ",처리 (處理),處理,,,,handling,,k,e,,,\n", out.String())
}

func setupTestFormat() (<-chan *pkg.Item, *bytes.Buffer) {
	items := make(chan *pkg.Item, 1)
	item := &pkg.Item{Hangul: "처리", Hanja: "處理", Def: pkg.Translation{English: "handling"}, Examples: []pkg.Translation{{English: "e", Korean: "k"}}}
	items <- item
	close(items)
	return items, new(bytes.Buffer)
}

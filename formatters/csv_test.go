package formatters

import (
	"bytes"
	"context"
	"testing"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

func TestFormatCSV(t *testing.T) {
	items, out := setupTestFormat()
	err := FormatCSV(context.Background(), items, out, map[string]string{})
	assert.Nil(t, err)
	assert.Equal(t, ",,처리,處理,,handling,,,,,,k,e,,\n", out.String())
}

func setupTestFormat() (<-chan *pkg.Item, *bytes.Buffer) {
	items := make(chan *pkg.Item, 1)
	item := &pkg.Item{Hangul: "처리", Hanja: "處理", Def: pkg.Translation{English: "handling"}, Examples: []pkg.Translation{{English: "e", Korean: "k"}}}
	items <- item
	close(items)
	return items, new(bytes.Buffer)
}

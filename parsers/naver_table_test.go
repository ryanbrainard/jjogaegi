package parsers

import (
	"context"
	"os"
	"testing"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

func TestParseNaverTable(t *testing.T) {
	in, err := os.Open("../testing/fixtures/naver_table.html")
	assert.Nil(t, err)

	items := make(chan *pkg.Item, 100)
	err = ParseNaverTable(context.Background(), in, items, map[string]string{})
	assert.Nil(t, err)

	assert.Equal(t, &pkg.Item{Hangul: "화살표화살", Hanja: "標", Def: pkg.Translation{English: "arrow"}}, <-items)
	assert.Equal(t, &pkg.Item{Hangul: "나열", Hanja: "羅列", Def: pkg.Translation{English: "[동사] list, (formal) enumerate"}}, <-items)
	assert.Equal(t, &pkg.Item{Hangul: "처리", Hanja: "處理", Def: pkg.Translation{English: "(일·사건 등의) handling; (쓰레기·폭탄 등의) disposal; (데이터 등의) processing, handle, deal with, take care of; (처분·제거하다) dispose of; (데이터 등을) process"}}, <-items)
}

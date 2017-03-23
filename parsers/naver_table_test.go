package parsers

import (
	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParseNaverTable(t *testing.T) {
	in := strings.NewReader(`처리處理 나열羅列 화살표
1 (일·사건 등의) handling 1 [동사] list
1 arrow
`)
	items := make(chan *pkg.Item, 100)
	ParseNaverTable(in, items, map[string]string{})
	assert.Equal(t, &pkg.Item{Hangul: "처리", Hanja: "處理", Def: pkg.Translation{English: "(일·사건 등의) handling"}}, <-items)
	assert.Equal(t, &pkg.Item{Hangul: "나열", Hanja: "羅列", Def: pkg.Translation{English: "[동사] list"}}, <-items)
	assert.Equal(t, &pkg.Item{Hangul: "화살표", Hanja: "", Def: pkg.Translation{English: "arrow"}}, <-items)
}

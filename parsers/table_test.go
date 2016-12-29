package parsers

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"github.com/ryanbrainard/jjogaegi"
)

func TestParseTable(t *testing.T) {
	in := strings.NewReader(`처리處理 나열羅列 화살표
1 (일·사건 등의) handling 1 [동사] list
1 arrow
`)
	items := make(chan *jjogaegi.Item, 100)
	ParseTable(in, items)
	assert.Equal(t, &jjogaegi.Item{Term: "처리", SubTerm: "處理", Def: "(일·사건 등의) handling"}, <-items)
	assert.Equal(t, &jjogaegi.Item{Term: "나열", SubTerm: "羅列", Def: "[동사] list"}, <-items)
	assert.Equal(t, &jjogaegi.Item{Term: "화살표", SubTerm: "", Def: "arrow"}, <-items)
}

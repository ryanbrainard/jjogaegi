package parsers

import (
	"strings"
	"testing"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

func TestParseMemriseList(t *testing.T) {
	in := strings.NewReader("안녕 라이언 \n Hello, Ryan")
	items := make(chan *pkg.Item, 100)
	ParseMemriseList(in, items, map[string]string{})
	assert.Equal(t, &pkg.Item{Id: "안녕 라이언", Hangul: "안녕 라이언", Def: pkg.Translation{English: "Hello, Ryan"}}, <-items)
}

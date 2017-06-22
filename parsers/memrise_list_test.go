package parsers

import (
	"os"
	"testing"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

func TestParseMemriseList(t *testing.T) {
	in, err := os.Open("fixtures/memrise_list.txt")
	assert.Nil(t, err)
	items := make(chan *pkg.Item, 100)
	ParseMemriseList(in, items, map[string]string{})
	assert.Equal(t, &pkg.Item{Id: "남성", Hangul: "남성", Def: pkg.Translation{English: "a man (not 남자)"}}, <-items)
	assert.Equal(t, &pkg.Item{Id: "너희", Hangul: "너희", Def: pkg.Translation{English: "you, you guys"}}, <-items)
}

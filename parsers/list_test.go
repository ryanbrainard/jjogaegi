package parsers

import (
	"github.com/ryanbrainard/jjogaegi"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParseListWithBullet(t *testing.T) {
	in := strings.NewReader(" • 안녕 라이언 Hello, Ryan")
	items := make(chan *jjogaegi.Item, 100)
	ParseList(in, items)
	assert.Equal(t, &jjogaegi.Item{Term: "안녕 라이언", Def: "Hello, Ryan"}, <-items)
}

func TestParseListWithNumberAndColon(t *testing.T) {
	in := strings.NewReader("1. 안녕 라이언: Hello, Ryan")
	items := make(chan *jjogaegi.Item, 100)
	ParseList(in, items)
	assert.Equal(t, &jjogaegi.Item{Term: "안녕 라이언", Def: "Hello, Ryan"}, <-items)
}

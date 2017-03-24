package parsers

import (
	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var xmlTestItems = []*pkg.Item{
	{
		Id:            "krdict:kor:15392:단어",
		Hangul:        "가입하다",
		Hanja:         "加入하다",
		Pronunciation: "가이파다",
		Def: pkg.Translation{
			Korean: "단체에 들어가거나 상품 및 서비스를 받기 위해 계약을 하다.",
		},
		Antonym: "탈퇴하다",
		Examples: []pkg.Translation{
			{
				Korean: "동아리에 가입하다.",
			},
			{
				Korean: "보험에 가입하다.",
			},
		},
	},
	{
		Id:            "krdict:kor:15404:단어",
		Hangul:        "갇히다",
		Pronunciation: "가치다",
		Def: pkg.Translation{
			Korean: "어떤 공간이나 상황에서 나가지 못하게 되다.",
		},
		Examples: []pkg.Translation{
			{
				Korean: "갇힌 몸.",
			},
			{
				Korean: "감옥에 갇히다.",
			},
		},
	},
}

func TestParseKrDictOppositesXML(t *testing.T) {
	in, err := os.Open("fixtures/kr_dict_opposites.xml")
	assert.Nil(t, err)

	items := make(chan *pkg.Item, 100)
	ParseKrDictXML(in, items, map[string]string{})
	assert.Equal(t, xmlTestItems[0], <-items)
	assert.Equal(t, xmlTestItems[1], <-items)
}

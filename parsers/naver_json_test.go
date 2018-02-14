package parsers

import (
	"os"
	"testing"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

var testItem = &pkg.Item{
	ExternalID: "naver:KE_729520861170427890b0340413d2af5a",
	Hangul:     "처리",
	Hanja:      "處理 (2)",
	Def: pkg.Translation{
		English: "1. (일·사건 등의) handling; (쓰레기·폭탄 등의) disposal; (데이터 등의) processing, handle, deal with, take care of; (처분·제거하다) dispose of; (데이터 등을) process  2. (물리적, 화학적) treatment, treat",
	},
	Examples: []pkg.Translation{
		{
			Korean:  "정보처리",
			English: "data[information] processing",
		},
		{
			Korean:  "폭탄처리반",
			English: "a bomb disposal squad",
		},
		{
			Korean:  "핵 재처리 시설",
			English: "nuclear reprocessing facilities",
		},
		{
			Korean:  "시신을 알코올로 처리하다",
			English: "treat a body with alcohol",
		},
	},
}

func TestParseNaverJSON(t *testing.T) {
	in, err := os.Open("fixtures/naver_json_sample.json")
	assert.Nil(t, err)

	items := make(chan *pkg.Item, 100)
	ParseNaverJSON(in, items, map[string]string{})
	assert.Equal(t, testItem, <-items)
}

func TestParseNaverJSONCallback(t *testing.T) {
	in, err := os.Open("fixtures/naver_json_sample.callback.js")
	assert.Nil(t, err)

	items := make(chan *pkg.Item, 100)
	ParseNaverJSON(in, items, map[string]string{})
	assert.Equal(t, testItem, <-items)
}

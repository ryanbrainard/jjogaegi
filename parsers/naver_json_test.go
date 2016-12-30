package parsers

import (
	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseNaverJSON(t *testing.T) {
	in, err := os.Open("naver_json_sample.json")
	assert.Nil(t, err)

	items := make(chan *pkg.Item, 100)
	ParseNaverJSON(in, items)
	assert.Equal(t, &pkg.Item{Term: "처리", SubTerm: "處理", Def: "1. (일·사건 등의) handling; (쓰레기·폭탄 등의) disposal; (데이터 등의) processing, handle, deal with, take care of; (처분·제거하다) dispose of; (데이터 등을) process  2. (물리적, 화학적) treatment, treat"}, <-items)
	assert.Equal(t, &pkg.Item{Term: "나열", SubTerm: "羅列", Def: "[동사] list, (formal) enumerate"}, <-items)
}

func TestParseNaverJSONCallback(t *testing.T) {
	in, err := os.Open("naver_json_sample.callback.js")
	assert.Nil(t, err)

	items := make(chan *pkg.Item, 100)
	ParseNaverJSON(in, items)
	assert.Equal(t, &pkg.Item{Term: "처리", SubTerm: "處理", Def: "1. (일·사건 등의) handling; (쓰레기·폭탄 등의) disposal; (데이터 등의) processing, handle, deal with, take care of; (처분·제거하다) dispose of; (데이터 등을) process  2. (물리적, 화학적) treatment, treat"}, <-items)
	assert.Equal(t, &pkg.Item{Term: "나열", SubTerm: "羅列", Def: "[동사] list, (formal) enumerate"}, <-items)
}

package parsers

import (
	"context"
	"os"
	"testing"

	"go.ryanbrainard.com/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

func TestParseNaverWordbook(t *testing.T) {
	in, err := os.Open("../testing/fixtures/naver_wordbook.html")
	assert.Nil(t, err)

	items := make(chan *pkg.Item, 100)
	err = ParseNaverWordbook(context.Background(), in, items, map[string]string{})
	assert.Nil(t, err)

	assert.Equal(t, &pkg.Item{
		Hangul: "화살표화살",
		Hanja:  "標",
		Def: pkg.Translation{
			English: "arrow",
		},
		Examples: []pkg.Translation{
			{
				Korean:  "화살표 방향으로 가시오",
				English: "Please follow this arrow.",
			},
		},
	}, <-items)

	assert.Equal(t, &pkg.Item{
		Hangul: "나열",
		Hanja:  "羅列",
		Def: pkg.Translation{
			English: "[동사] list, (formal) enumerate",
		},
		Examples: []pkg.Translation{
			{
				Korean:  "단편적인 정보를 나열하다",
				English: "list fragments of information",
			},
			{
				Korean:  "단편적인 정보를 나열하다",
				English: "enumerate bits of information",
			},
		},
	}, <-items)

	assert.Equal(t, &pkg.Item{
		Hangul: "처리",
		Hanja:  "處理",
		Def: pkg.Translation{
			English: "(일·사건 등의) handling; (쓰레기·폭탄 등의) disposal; (데이터 등의) processing, handle, deal with, take care of; (처분·제거하다) dispose of; (데이터 등을) process",
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
		},
	}, <-items)
}

package parsers

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ryanbrainard.com/jjogaegi/pkg"
)

func TestParseNaverWordbookJSON(t *testing.T) {
	in, err := os.Open("../testing/fixtures/naver_wordbook.json")
	require.Nil(t, err)

	items := make(chan *pkg.Item, 100)
	err = ParseNaverWordbookJSON(context.Background(), in, items, map[string]string{pkg.OPT_DEBUG: "true"})
	assert.Nil(t, err)

	assert.Equal(t, &pkg.Item{
		ExternalID:    "b2f004c463a04303a917fa24bef83906",
		Pronunciation: "[nayeol]",
		Hangul:        "나열",
		Hanja:         "羅列",
		Def: pkg.Translation{
			English: "[동사] list, (formal) enumerate",
		},
		Examples: []pkg.Translation{
			{
				Korean:  "단편적인 정보를 나열하다",
				English: "list fragments of information",
			},
		},
	}, <-items)

	assert.Equal(t, &pkg.Item{
		ExternalID:    "cea6a785192d44608c4e6207442ef68d",
		Pronunciation: "[hwasalpyo]",
		Hangul:        "화살표",
		Hanja:         "",
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
		ExternalID:    "729520861170427890b0340413d2af5a",
		Pronunciation: "[cheo-ri]",
		Hangul:        "처리",
		Hanja:         "處理",
		Def: pkg.Translation{
			English: "(일·사건 등의) handling; (쓰레기·폭탄 등의) disposal; (데이터 등의) processing, handle, deal with, take care of; (처분·제거하다) dispose of; (데이터 등을) process",
		},
		Examples: []pkg.Translation{
			{
				Korean:  "정보처리",
				English: "data[information] processing",
			},
		},
	}, <-items)
}

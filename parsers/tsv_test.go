package parsers

import (
	"os"
	"testing"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

func TestParseTSV(t *testing.T) {
	in, err := os.Open("../testing/fixtures/anki.tsv")
	assert.Nil(t, err)

	items := make(chan *pkg.Item, 100)
	err = ParseTSV(in, items, map[string]string{})
	assert.Nil(t, err)

	expected := &pkg.Item{
		NoteID:        "1517799325880",
		ExternalID:    "krdict:kor:50010:단어",
		Hangul:        "막상막하",
		Hanja:         "莫上莫下",
		Pronunciation: "막쌍마카",
		AudioTag:      "[sound:makssangmaka.wav]",
		ImageTag:      "<img src=\"paste-2512555868520.jpg\" />",
		Def: pkg.Translation{
			Korean:  "누가 더 나은지 가릴 수 없을 만큼 차이가 거의 없음.",
			English: ":=",
		},
		Antonym: "",
		Examples: []pkg.Translation{
			pkg.Translation{
				Korean:  "막상막하의 경기.",
				English: "",
			},
			pkg.Translation{
				Korean:  "이번 의원 선거에서는 여당과 야당 후보의 득표율이 막상막하였다.",
				English: "",
			},
		},
	}

	assert.Equal(t, expected, <-items)

}

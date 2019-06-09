package parsers

import (
	"context"
	"os"
	"testing"

	"ryanbrainard.com/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

var xmlTestItems = []*pkg.Item{
	{
		ExternalID:    "krdict:kor:15392:단어",
		Hangul:        "가입하다",
		Hanja:         "加入하다",
		Pronunciation: "가이파다",
		AudioTag:      "http://dicmedia.korean.go.kr:8899/multimedia/sound_file/giyeok_2006/gaipada.wav",
		Def: pkg.Translation{
			Korean: "단체에 들어가거나 상품 및 서비스를 받기 위해 계약을 하다.",
			// English: "join; sign up := To join a group or sign up for goods and services.",
		},
		Antonym: "탈퇴하다",
		Examples: []pkg.Translation{
			{
				Korean: "동아리에 가입하다.",
			},
		},
	},
	{
		ExternalID:    "krdict:kor:82614:단어",
		Hangul:        "탈퇴하다",
		Hanja:         "脫退하다",
		Pronunciation: "탈퇴하다",
		AudioTag:      "http://dicmedia.korean.go.kr:8899/multimedia/sound_file/tieut_2005/taltoehada.wav",
		Def: pkg.Translation{
			Korean: "소속해 있던 조직이나 단체에서 관계를 끊고 나오다.",
			// English: "withdraw; drop out; leave := To end one's relationship with an organization or group one had belonged to and leave.",
		},
		Antonym: "가입하다",
		Examples: []pkg.Translation{
			{
				Korean: "신속히 탈퇴하다.",
			},
			{
				Korean: "나는 공부에 전념하기 위해 동아리에서 탈퇴했다.",
			},
		},
	},
}

func TestParseKrDictOppositesXML(t *testing.T) {
	in, err := os.Open("../testing/fixtures/kr_dict_opposites.xml")
	assert.Nil(t, err)

	items := make(chan *pkg.Item, 100)
	err = ParseKrDictXML(context.Background(), in, items, map[string]string{pkg.OPT_MEDIADIR: "/tmp"})
	assert.Nil(t, err)
	assert.Equal(t, xmlTestItems[0], <-items)
	assert.Equal(t, xmlTestItems[1], <-items)
}

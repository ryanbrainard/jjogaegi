package interceptors

import (
	"os"
	"testing"

	"launchpad.net/xmlpath"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

var xmlTestItems = []*pkg.Item{
	{
		ExternalID:    "krdict:kor:15392:단어",
		Hangul:        "가입하다",
		Hanja:         "加入",
		Pronunciation: "가이파다",
		// AudioTag:      "[sound:gaipada.wav]",
		Def: pkg.Translation{
			Korean:  "단체에 들어가거나 상품 및 서비스를 받기 위해 계약을 하다.",
			English: "join; sign up := To join a group or sign up for goods and services.",
		},
		Antonym: "탈퇴하다",
		Examples: []pkg.Translation{
			{
				Korean: "동아리에 가입하다.",
			},
			{
				Korean: "국제 연합에 가입한 회원국은 총 백구십이 개국이다.",
			},
		},
	},
	{
		ExternalID:    "krdict:kor:82614:단어",
		Hangul:        "탈퇴하다",
		Hanja:         "脫退",
		Pronunciation: "탈퇴하다",
		// AudioTag:      "[sound:taltoehada.wav]",
		Def: pkg.Translation{
			Korean:  "소속해 있던 조직이나 단체에서 관계를 끊고 나오다.",
			English: "withdraw; drop out; leave := To end one's relationship with an organization or group one had belonged to and leave.",
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

func TestKrDictEnhance(t *testing.T) {
	if os.Getenv("KRDICT_API_KEY") == "" {
		t.Skip("KRDICT_API_KEY not set")
	}

	for _, expectedItem := range xmlTestItems {
		t.Run(expectedItem.Hangul, func(tr *testing.T) {
			actualItem := &pkg.Item{
				ExternalID: expectedItem.ExternalID,
			}

			err := KrDictEnhance(actualItem, map[string]string{})
			assert.NoError(t, err)
			assert.Equal(tr, expectedItem, actualItem)
		})
	}
}
func TestGetEnglishDefinition(t *testing.T) {
	in, err := os.Open("../parsers/fixtures/kr_dict_en_15392.xml")
	assert.Nil(t, err)
	node, err := xmlpath.Parse(in)
	assert.Nil(t, err)

	assert.Equal(t, "join; sign up := To join a group or sign up for goods and services.", getEnglishDefinition(node))
}
func TestGetWordGrade(t *testing.T) {
	in, err := os.Open("../parsers/fixtures/kr_dict_en_15392-mod.xml")
	assert.Nil(t, err)
	node, err := xmlpath.Parse(in)
	assert.Nil(t, err)

	assert.Equal(t, "고급", getWordGrade(node))
}

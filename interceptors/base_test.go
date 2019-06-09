package interceptors

import (
	"ryanbrainard.com/jjogaegi/pkg"
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
	{
		ExternalID:    "krdict:kor:71793:단어",
		Hangul:        "참기름",
		Pronunciation: "참기름",
		Def: pkg.Translation{
			Korean:  "참깨로 짠 기름.",
			English: "sesame oil := The oil extracted from sesame seeds.",
		},
		Examples: []pkg.Translation{
			{
				Korean: "고소한 참기름.",
			},
			{
				Korean: "이 참기름은 방앗간에서 직접 참깨로 짠 것이라서 그런지 맛과 향이 다른 것보다 더 고소하다.",
			},
		},
		Grade:    "중급",
		ImageTag: "http://dicmedia.korean.go.kr:8899/front/search/searchResultView.do?file_no=197671",
	},
}

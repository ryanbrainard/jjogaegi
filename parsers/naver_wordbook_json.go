package parsers

import (
	"context"
	"encoding/json"
	"go.ryanbrainard.com/jjogaegi/pkg"
	"io"
	"log"
)

func ParseNaverWordbookJSON(ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	page := &NaverWordbookPage{}

	if err := json.NewDecoder(r).Decode(page); err != nil {
		return err
	}

	for _, mitem := range page.Data.MItems {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		item := &pkg.Item{
			ExternalID: mitem.EntryID,
		}

		var content interface{}
		var handleContent func(*pkg.Item, interface{})

		switch mitem.ContentType {
		case 0:
			content = &NaverWordbookItemContentType0{}
			handleContent = func(item *pkg.Item, content interface{}) {
				handleContentType0(item, content.(*NaverWordbookItemContentType0))
			}
		case 1:
			content = &NaverWordbookItemContentType1{}
			handleContent = func(item *pkg.Item, content interface{}) {
				handleContentType1(item, content.(*NaverWordbookItemContentType1))
			}
		default:
			log.Printf("unknown content type=%d name=%s", mitem.ContentType, mitem.Name)
			continue
		}

		if err := json.Unmarshal([]byte(mitem.Content), &content); err != nil {
			return err
		}

		pkg.Debug(options, "fn=ParseNaverWordbookJSON entry_id=%s content=%+v", mitem.EntryID, content)

		handleContent(item, content)

		pkg.Debug(options, "fn=ParseNaverWordbookJSON entry_id=%s item=%+v", mitem.EntryID, item)

		items <- item
	}

	return nil
}

func handleContentType0(item *pkg.Item, content *NaverWordbookItemContentType0) {
	if len(content.Entry.Members) > 0 {
		member := content.Entry.Members[0]

		item.Hangul = member.EntryName
		item.Hanja = member.OriginLanguage

		if len(member.Prons) > 0 {
			item.Pronunciation = member.Prons[0].PronSymbol
		}
	}

	if len(content.Entry.Means) > 0 {
		mean := content.Entry.Means[0]

		item.Def = pkg.Translation{
			English: mean.ShowMean,
		}

		for _, example := range mean.Examples {
			handleExample(item, example)
		}
	}
}

func handleContentType1(item *pkg.Item, content *NaverWordbookItemContentType1) {
	handleExample(item, content.Example)
}

func handleExample(item *pkg.Item, example NaverWordbookExample) {
	var englishExample string
	var koreanExample string

	switch example.Language {
	case "en":
		englishExample = example.ShowExample
		koreanExample = findTranslations(example.Translations, "ko")
	case "ko":
		koreanExample = example.ShowExample
		englishExample = findTranslations(example.Translations, "en")
	}

	item.Examples = append(item.Examples, pkg.Translation{
		Korean:  koreanExample,
		English: englishExample,
	})
}

func findTranslations(translations []NaverWordbookTranslation, language string) string {
	for _, translation := range translations {
		if translation.Language == language {
			return translation.ShowTranslation
		}
	}
	return ""
}

type NaverWordbookPage struct {
	Meta struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"meta"`
	Data struct {
		OverLastPage bool   `json:"over_last_page"`
		NextCursor   string `json:"next_cursor"`
		MTotal       int    `json:"m_total"`
		MItems       []struct {
			ID          string      `json:"id"`
			EntryID     string      `json:"entryId"`
			WordbookID  string      `json:"wordbookId"`
			Service     string      `json:"service"`
			DicType     string      `json:"dicType"`
			Name        string      `json:"name"`
			ContentType int         `json:"contentType"`
			SaveType    string      `json:"saveType"`
			SourceType  string      `json:"sourceType"`
			Content     string      `json:"content"`
			Memo        interface{} `json:"memo"`
			CreateDate  int64       `json:"createDate"`
			Read        bool        `json:"read"`
		} `json:"m_items"`
	} `json:"data"`
}

type NaverWordbookItemContentType0 struct {
	Entry struct {
		EntryID           string `json:"entry_id"`
		Service           string `json:"service"`
		DictType          string `json:"dict_type"`
		DictCid           string `json:"dict_cid"`
		MixPron           string `json:"mix_pron"`
		KoreanMixPron     string `json:"korean_mix_pron"`
		StrokeShowFile    string `json:"stroke_show_file"`
		StrokeCnt         string `json:"stroke_cnt"`
		RadicalID         string `json:"radical_id"`
		ShowRadical       string `json:"show_radical"`
		KoreanRadicalName string `json:"korean_radical_name"`
		SourceType        string `json:"source_type"`
		SourceName        string `json:"source_name"`
		Language          string `json:"language"`
		Members           []struct {
			MemberID        string `json:"member_id"`
			EntryName       string `json:"entry_name"`
			SuperScript     string `json:"super_script"`
			EntryImportance string `json:"entry_importance"`
			OriginLanguage  string `json:"origin_language"`
			Kanji           string `json:"kanji"`
			Prons           []struct {
				PronID           string `json:"pron_id"`
				PronSymbol       string `json:"pron_symbol"`
				PronAlias        string `json:"pron_alias"`
				PronType         string `json:"pron_type"`
				KoreanPronSymbol string `json:"korean_pron_symbol"`
				OrderSeq         int    `json:"order_seq"`
				MalePronFile     string `json:"male_pron_file"`
				FemalePronFile   string `json:"female_pron_file"`
			} `json:"prons"`
		} `json:"members"`
		Means []struct {
			MeanID          string `json:"mean_id"`
			PartID          string `json:"part_id"`
			PartName        string `json:"part_name"`
			ShowMean        string `json:"show_mean"`
			Language        string `json:"language"`
			DescriptionJSON struct {
			} `json:"description_json"`
			MeanType  string                 `json:"mean_type"`
			MeanLevel int                    `json:"mean_level"`
			MeanSeq   interface{}            `json:"mean_seq"`
			Examples  []NaverWordbookExample `json:"examples"`
		} `json:"means"`
	} `json:"entry"`
}

type NaverWordbookItemContentType1 struct {
	Example NaverWordbookExample `json:"examples"`
}

type NaverWordbookExample struct {
	ExampleID       string                     `json:"example_id"`
	ShowExample     string                     `json:"show_example"`
	Language        string                     `json:"language"`
	ExampleType     string                     `json:"example_type"`
	DictCid         string                     `json:"dict_cid"`
	DictType        string                     `json:"dict_type"`
	ExampleSource   string                     `json:"example_source"`
	ExamplePron     string                     `json:"example_pron"`
	ExamplePronFile string                     `json:"example_pron_file"`
	Translations    []NaverWordbookTranslation `json:"translations"`
}

type NaverWordbookTranslation struct {
	Language        string `json:"language"`
	TranslationID   string `json:"translation_id"`
	ShowTranslation string `json:"show_translation"`
}

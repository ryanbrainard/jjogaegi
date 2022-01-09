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

		if mitem.ContentType != 0 {
			log.Printf("unknown content type=%d name=%s", mitem.ContentType, mitem.Name)
			continue
		}

		content := &NaverWordbookItemContent{}
		if err := json.Unmarshal([]byte(mitem.Content), content); err != nil {
			return err
		}

		pkg.Debug(options, "fn=ParseNaverWordbookJSON entry_id=%s content=%+v", content.Entry.EntryID, content)

		item := &pkg.Item{
			ExternalID: content.Entry.EntryID,
		}

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
				var english string
				if len(example.Translations) > 0 && example.Translations[0].Language == "en" {
					english = example.Translations[0].ShowTranslation
				}

				item.Examples = append(item.Examples, pkg.Translation{
					Korean:  example.ShowExample,
					English: english,
				})
			}
		}

		pkg.Debug(options, "fn=ParseNaverWordbookJSON entry_id=%s item=%+v", content.Entry.EntryID, item)

		items <- item
	}

	return nil
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

type NaverWordbookItemContent struct {
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
			MeanType  string      `json:"mean_type"`
			MeanLevel int         `json:"mean_level"`
			MeanSeq   interface{} `json:"mean_seq"`
			Examples  []struct {
				ExampleID       string `json:"example_id"`
				ShowExample     string `json:"show_example"`
				Language        string `json:"language"`
				ExampleType     string `json:"example_type"`
				DictCid         string `json:"dict_cid"`
				ExamplePron     string `json:"example_pron"`
				ExamplePronFile string `json:"example_pron_file"`
				Translations    []struct {
					Language        string `json:"language"`
					TranslationID   string `json:"translation_id"`
					ShowTranslation string `json:"show_translation"`
				} `json:"translations"`
			} `json:"examples"`
		} `json:"means"`
	} `json:"entry"`
}

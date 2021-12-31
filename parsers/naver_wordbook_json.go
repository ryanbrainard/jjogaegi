package parsers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.ryanbrainard.com/jjogaegi/pkg"
	"io"
)

func ParseNaverWordbookJSON(ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	page := &NaverWordbookPage{}

	if err := json.NewDecoder(r).Decode(page); err != nil {
		return err
	}

	for _, mitem := range page.Data.MItems {
		if mitem.ContentType != 0 {
			return fmt.Errorf("unknown content type: %d", mitem.ContentType)
		}

		content := &NaverWordbookItemContent{}
		if err := json.Unmarshal([]byte(mitem.Content), content); err != nil {
			return err
		}

		// using first element of each. TOOD: filter
		member := content.Entry.Members[0]
		prons := member.Prons[0]
		mean := content.Entry.Means[0]
		example := mean.Examples[0]
		translation := example.Translations[0]

		items <- &pkg.Item{
			ExternalID: content.Entry.EntryID,
			Hangul:     member.EntryName,
			Hanja:      member.OriginLanguage,
			Def: pkg.Translation{
				English: mean.ShowMean,
			},
			Pronunciation: prons.PronSymbol,
			Examples: []pkg.Translation{
				{
					Korean:  example.ShowExample,
					English: translation.ShowTranslation,
				},
			},
		}
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
			MeanType  string `json:"mean_type"`
			MeanLevel int    `json:"mean_level"`
			MeanSeq   string `json:"mean_seq"`
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

package interceptors

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"launchpad.net/xmlpath"
)

func KrDictEnhance(item *pkg.Item, options map[string]string) error {
	// format: strings.Join([]string{"krdict", lang, entryId, lexicalUnit}, ":")
	idSplit := strings.Split(item.ExternalID, ":")
	if len(idSplit) != 4 || idSplit[0] != "krdict" {
		return nil
	}
	entryID := idSplit[2]

	entry, err := fetchEntryNode(entryID)
	if err != nil {
		return err
	}

	if item.Hangul == "" {
		item.Hangul = get(entry, "/channel/item/word_info/word")
	}

	if item.Hanja == "" {
		item.Hanja = get(entry, "/channel/item/word_info/original_language_info[language_type='한자']/original_language")
	}

	if item.Pronunciation == "" {
		item.Pronunciation = get(entry, "/channel/item/word_info/pronunciation_info/pronunciation")
	}

	// TODO: broken because missing
	// if item.AudioTag == "" {
	// }

	if item.Def.Korean == "" {
		item.Def.Korean = get(entry, "/channel/item/word_info/sense_info/definition")
	}

	if item.Def.English == "" {
		item.Def.English = getEnglishDefinition(entry)
	}

	if item.Antonym == "" {
		item.Antonym = get(entry, "/channel/item/word_info/sense_info/rel_info[type='반대말']/word")
	}

	if item.Examples == nil {
		item.Examples = []pkg.Translation{}
	}

	if len(item.Examples) == 0 {
		if example := getExample(entry, "구"); example != nil {
			item.Examples = append(item.Examples, *example)
		}
	}

	if len(item.Examples) == 1 {
		if example := getExample(entry, "문장"); example != nil {
			item.Examples = append(item.Examples, *example)
		}
	}

	if item.ImageTag == "" {
		// TODO: why isn't filter working?
		// get(entry, "/channel/item/word_info/sense_info/multimedia_info[type='사진']/link")
		item.ImageTag = get(entry, "/channel/item/word_info/sense_info/multimedia_info/link")
	}

	if item.Grade == "" || item.Grade == "없음" {
		item.Grade = getWordGrade(entry)
	}

	return nil
}

func fetchEntryNode(entryID string) (*xmlpath.Node, error) {
	if os.Getenv("KRDICT_API_KEY") == "" {
		panic("KRDICT_API_KEY not set.")
	}

	url := fmt.Sprintf("https://krdict.korean.go.kr/api/view?key=%s&type_search=view&method=TARGET_CODE&part=word&q=%s&sort=dict&translated=y&trans_lang=1", os.Getenv("KRDICT_API_KEY"), entryID)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("download type=eng url=%q err=%q", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	return xmlpath.Parse(resp.Body)
}

func getEnglishDefinition(node *xmlpath.Node) string {
	transPath := "/channel/item/word_info/sense_info/translation"
	transWord := get(node, transPath+"/trans_word")
	transDfn := get(node, transPath+"/trans_dfn")

	return transWord + " := " + transDfn
}

func getWordGrade(node *xmlpath.Node) string {
	grade := get(node, "/channel/item/word_info/word_grade")
	switch grade {
	case "없음":
		return ""
	default:
		// TODO: consider changing to numbers
		return grade
	}
}

func getExample(node *xmlpath.Node, exampleType string) *pkg.Translation {
	examplesIter := xmlpath.MustCompile("/channel/item/word_info/sense_info/example_info").Iter(node)
	for {
		if !examplesIter.Next() {
			break
		}

		exampleNode := examplesIter.Node()

		if get(exampleNode, "type") == exampleType {
			return &pkg.Translation{Korean: get(exampleNode, "example")}
		}
	}
	return nil
}

func get(node *xmlpath.Node, xpath string) string {
	path := xmlpath.MustCompile(xpath)

	if value, ok := path.String(node); ok {
		return strings.TrimSpace(value)
	}

	return ""
}

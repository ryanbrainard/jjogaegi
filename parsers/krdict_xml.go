package parsers

import (
	"io"
	"log"
	"strings"

	"net/http"

	"fmt"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"launchpad.net/xmlpath"
	"os"
)

const supportedLanguage = "kor"
const supportedLexicalUnit = "단어"

func ParseKrDictXML(r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	rootNode, err := xmlpath.Parse(r)
	if err != nil {
		return err
	}

	lang := get(rootNode, "/LexicalResource/Lexicon/feat[@att='language']/@val")
	if lang != supportedLanguage {
		return fmt.Errorf("Only %q supported.", supportedLanguage)
	}

	entriesIter := xmlpath.MustCompile("/LexicalResource/Lexicon/LexicalEntry").Iter(rootNode)
	for {
		if !entriesIter.Next() {
			break
		}
		entryNode := entriesIter.Node()

		lexicalUnit := get(entryNode, "feat[@att='lexicalUnit']/@val")
		if lexicalUnit != supportedLexicalUnit {
			continue
		}

		entryId := get(entryNode, ".[@att='id']/@val")

		item := &pkg.Item{
			Id:            strings.Join([]string{"krdict", lang, entryId, lexicalUnit}, ":"),
			Hangul:        get(entryNode, "Lemma/feat/@val"),
			Hanja:         get(entryNode, "feat[@att='origin']/@val"),
			Pronunciation: get(entryNode, "WordForm/feat[@att='pronunciation']/@val"),
			AudioURL:      get(entryNode, "WordForm/feat[@att='sound']/@val"),
			ImageURL:      get(entryNode, "Sense/Multimedia/feat[@att='url']/@val"),
			Antonym:       get(entryNode, "Sense/SenseRelation/feat[@val='반대말']/../feat[@att='lemma']/@val"),
			Def: pkg.Translation{
				Korean:  get(entryNode, "Sense/feat[@att='definition']/@val"),
				English: fetchEnglishDefinition(entryId),
			},
		}

		exampleType := ""
		examplesIter := xmlpath.MustCompile("Sense/SenseExample").Iter(entryNode)
		for {
			if !examplesIter.Next() {
				break
			}
			exampleNode := examplesIter.Node()

			nextExampleType := get(exampleNode, "feat[@att='type']/@val")
			if nextExampleType == exampleType {
				continue // choose one of each example type
			}
			if nextExampleType == "대화" {
				continue
			}
			exampleType = nextExampleType

			item.Examples = append(item.Examples, pkg.Translation{
				Korean: get(exampleNode, "feat[@att='example']/@val"),
			})
		}

		items <- item
	}

	return nil
}

func get(node *xmlpath.Node, xpath string) string {
	path := xmlpath.MustCompile(xpath)

	if value, ok := path.String(node); ok {
		return value
	}

	return ""
}

func fetchEnglishDefinition(entryId string) string {
	url := fmt.Sprintf("https://krdict.korean.go.kr/api/view?key=%s&type_search=view&method=TARGET_CODE&part=word&q=%s&sort=dict&translated=y&trans_lang=1", os.Getenv("KRDICT_API_KEY"), entryId)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("download type=eng url=%q err=%q", url, err)
		return ""
	}
	defer resp.Body.Close()

	return extractEnglishDefinition(resp.Body)
}

func extractEnglishDefinition(r io.Reader) string {
	node, err := xmlpath.Parse(r)
	if err != nil {
		return ""
	}

	transPath := "/channel/item/word_info/sense_info/translation"
	transWord := get(node, transPath + "/trans_word")
	transDfn := get(node, transPath + "/trans_dfn")

	return transWord + " := " + transDfn
}

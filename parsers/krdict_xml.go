package parsers

import (
	"github.com/ryanbrainard/jjogaegi/pkg"
	"io"
	"launchpad.net/xmlpath"
	"log"
	"strings"
)

const SupportedLanguage = "kor"

func ParseKrDictXML(r io.Reader, items chan<- *pkg.Item, options map[string]string) {
	rootNode, err := xmlpath.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	lang := get(rootNode, "/LexicalResource/Lexicon/feat[@att='language']/@val")
	if lang != SupportedLanguage {
		log.Fatalf("Only %q supported.", SupportedLanguage)
	}

	entriesIter := xmlpath.MustCompile("/LexicalResource/Lexicon/LexicalEntry").Iter(rootNode)
	for {
		if !entriesIter.Next() {
			break
		}
		entryNode := entriesIter.Node()

		item := &pkg.Item{
			Id:            strings.Join([]string{"krdict", lang, get(entryNode, ".[@att='id']/@val"), get(entryNode, "feat[@att='lexicalUnit']/@val")}, ":"),
			Hangul:        get(entryNode, "Lemma/feat/@val"),
			Hanja:         get(entryNode, "feat[@att='origin']/@val"),
			Pronunciation: get(entryNode, "WordForm/feat[@att='pronunciation']/@val"),
			AudioURL:      get(entryNode, "WordForm/feat[@att='sound']/@val"),
			Antonym:       get(entryNode, "Sense/SenseRelation/feat[@val='반대말']/../feat[@att='lemma']/@val"),
			Def: pkg.Translation{
				Korean: get(entryNode, "Sense/feat[@att='definition']/@val"),
			},
		}

		examplesIter := xmlpath.MustCompile("Sense/SenseExample").Iter(entryNode)
		for {
			if !examplesIter.Next() {
				break
			}
			exampleNode := examplesIter.Node()

			item.Examples = append(item.Examples, pkg.Translation{
				Korean: get(exampleNode, "feat[@att='example']/@val"),
			})
		}

		items <- item
	}

	close(items)
}

func get(node *xmlpath.Node, xpath string) string {
	path := xmlpath.MustCompile(xpath)

	//log.Printf("xpath=%q exists=%t", xpath, path.Exists(node))

	if value, ok := path.String(node); ok {
		return value
	}

	return ""
}

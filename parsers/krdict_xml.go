package parsers

import (
	"io"
	"log"
	"strings"

	"net/http"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"golang.org/x/net/html"
	"launchpad.net/xmlpath"
)

const supportedLanguage = "kor"
const supportedLexicalUnit = "단어"

func ParseKrDictXML(r io.Reader, items chan<- *pkg.Item, options map[string]string) {
	rootNode, err := xmlpath.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	lang := get(rootNode, "/LexicalResource/Lexicon/feat[@att='language']/@val")
	if lang != supportedLanguage {
		log.Fatalf("Only %q supported.", supportedLanguage)
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

	close(items)
}

func get(node *xmlpath.Node, xpath string) string {
	path := xmlpath.MustCompile(xpath)

	if value, ok := path.String(node); ok {
		return value
	}

	return ""
}

func fetchEnglishDefinition(entryId string) string {
	url := "https://krdict.korean.go.kr/eng/dicSearch/SearchView?nation=eng&ParaWordNo=" + entryId

	log.Printf("download type=eng url=%q", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("download type=eng url=%q err=%q", url, err)
		return ""
	}
	defer resp.Body.Close()

	return extractEnglishDefinition(resp.Body)
}

func extractEnglishDefinition(r io.Reader) string {
	z := html.NewTokenizer(r)
	out := ""
	stack := &tokenStack{}
	processed := 0
	for {
		switch z.Next() {
		case html.StartTagToken:
			if string(z.Raw()) == `<p class="theme1 multiTrans manyLang6 transFont6" style="margin-bottom: 13px;">` {
				stack.Push(z.Token())
			}
			if stack.Depth() > 0 && string(z.Raw()) == `<strong>` {
				stack.Push(z.Token())
			}
			if string(z.Raw()) == `<p class="sub_p1 manyLang6 multiSenseDef defFont6" style="margin-left: 20px;line-height: 20px;">` {
				stack.Push(z.Token())
				out += " := "
			}
		case html.TextToken:
			if stack.Depth() == 1 {
				out += strings.Trim(string(z.Text()), " \n\t")
			}
		case html.EndTagToken:
			if stack.Depth() > 0 && stack.Peek().Data == z.Token().Data {
				stack.Pop()
				processed++
				if processed > 2 {
					return out
				}
			}
		case html.ErrorToken:
			return out
		}
	}
	return out
}

type tokenStack struct {
	stack []html.Token
}

func (ts *tokenStack) Push(v html.Token) {
	ts.stack = append(ts.stack, v)
}

func (ts *tokenStack) Pop() html.Token {
	res := ts.stack[ts.Depth()-1]
	ts.stack = ts.stack[:ts.Depth()-1]
	return res
}

func (ts *tokenStack) Peek() html.Token {
	return ts.stack[ts.Depth()-1]
}

func (ts *tokenStack) Depth() int {
	return len(ts.stack)
}

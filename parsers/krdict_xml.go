package parsers

import (
	"github.com/ryanbrainard/jjogaegi/pkg"
	"io"
	"launchpad.net/xmlpath"
	"log"
)

func ParseKrDictXML(r io.Reader, items chan<- *pkg.Item, options map[string]string) {
	rootNode, err := xmlpath.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	entriesIter := xmlpath.MustCompile("/LexicalResource/Lexicon/LexicalEntry").Iter(rootNode)
	for {
		if !entriesIter.Next() {
			break
		}
		entryNode := entriesIter.Node()

		item := &pkg.Item{
			Hangul:   get(entryNode, "Lemma/feat/@val"),
			Hanja:    get(entryNode, "feat[@att='origin']/@val"),
			Def:      get(entryNode, "Sense/feat[@att='definition']/@val"),
		}

		examplesIter := xmlpath.MustCompile("Sense/SenseExample").Iter(entryNode)
		for {
			if !examplesIter.Next() {
				break
			}
			exampleNode := examplesIter.Node()

			item.Examples = append(item.Examples, pkg.Example{
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

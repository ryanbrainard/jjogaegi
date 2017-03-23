package parsers

import (
	"github.com/ryanbrainard/jjogaegi/pkg"
	"io"
	"launchpad.net/xmlpath"
	"log"
)

func ParseKrDictXML(r io.Reader, items chan<- *pkg.Item, options map[string]string) {
	root, err := xmlpath.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	iter := xmlpath.MustCompile("/LexicalResource/Lexicon/LexicalEntry").Iter(root)

	for {
		if !iter.Next() {
			break
		}

		node := iter.Node()

		items <- &pkg.Item{
			Hangul: get(node, "Lemma/feat/@val"),
			Hanja:  get(node, "feat[@att='origin']/@val"),
			Def:    get(node, "Sense/feat[@att='definition']/@val"),
			Examples: []pkg.Example{
				{
					Korean: get(node, "Sense/SenseExample/feat[@att='example']/@val"),
				},
			},
		}
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

package parsers

import (
	"context"
	"ryanbrainard.com/jjogaegi/pkg"
	"io"
	"launchpad.net/xmlpath"
	"strings"
)

func ParseNaverWordbook(ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	rootNode, err := xmlpath.ParseHTML(r)
	if err != nil {
		return err
	}

	entriesIter := xmlpath.MustCompile("//div[@class='print_article']").Iter(rootNode)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if !entriesIter.Next() {
			break
		}
		entryNode := entriesIter.Node()

		hangulHanja := pkg.XpathString(entryNode, "h2[@class='entry']/em[@class='word']")
		hangul, hanja := splitHangul(hangulHanja)

		meaningIter := xmlpath.MustCompile("ol[@class='mean']/li[@class='row']").Iter(entryNode)
		if !meaningIter.Next() {
			continue
		}
		meaningNode := meaningIter.Node() // only get first meaning, if it exists

		item := &pkg.Item{
			Hangul: strings.TrimSpace(hangul),
			Hanja:  strings.TrimSpace(hanja),
			Def: pkg.Translation{
				English: pkg.XpathString(meaningNode, "p[@class='speech']"),
			},
		}

		examplesIter := xmlpath.MustCompile("p[@class='ex']").Iter(meaningNode)
		for {
			if !examplesIter.Next() {
				break
			}
			exampleNode := examplesIter.Node()

			english, korean := splitHangulReverse(exampleNode.String())

			item.Examples = append(item.Examples, pkg.Translation{
				English: english,
				Korean:  korean,
			})
		}

		items <- item
	}

	return nil
}

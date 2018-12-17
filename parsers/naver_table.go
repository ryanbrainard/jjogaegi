package parsers

import (
	"context"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"io"
	"launchpad.net/xmlpath"
	"strings"
)

func ParseNaverTable(ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
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

		item := &pkg.Item{
			Hangul: strings.TrimSpace(hangul),
			Hanja:  strings.TrimSpace(hanja),
			//Def: pkg.Translation{
			//	English: pkg.XpathString(entryNode, "//p[@class='speech']"),
			//},
		}

		items <- item
	}

	return nil
}

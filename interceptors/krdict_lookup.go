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

func KrDictLookup(item *pkg.Item, options map[string]string) error {
	if options[pkg.OPT_KRDICT_LOOKUP] != "true" {
		return nil
	}

	if strings.HasPrefix(item.ExternalID, "krdict") {
		return nil
	}

	q := item.Hangul
	if item.Hanja != "" {
		q = item.Hanja
	}

	results, err := search(q)
	if err != nil {
		return err
	}

	resultsIntr := xmlpath.MustCompile("/channel/item").Iter(results)
	choices := []*xmlpath.Node{}
	for {
		if !resultsIntr.Next() {
			break
		}

		result := resultsIntr.Node()

		if item.Hangul != pkg.XpathString(result, "word") {
			continue
		}

		choices = append(choices, result)
	}

	switch len(choices) {
	case 0:
		println("no results: ", item.Hangul)
	case 1:
		item.ExternalID = pkg.KrDictID("kor", pkg.XpathString(choices[0], "target_code"), "단어")
	default:
		println("multiple results: ", item.Hangul)
	}

	return nil
}

func search(q string) (*xmlpath.Node, error) {
	if os.Getenv("KRDICT_API_KEY") == "" {
		panic("KRDICT_API_KEY not set.")
	}

	url := fmt.Sprintf("https://krdict.korean.go.kr/api/search?key=%s&type_search=search&part=word&q=%s&sort=dict&translated=y&trans_lang=1", os.Getenv("KRDICT_API_KEY"), q)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("search type=eng url=%q err=%q", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	return xmlpath.Parse(resp.Body)
}

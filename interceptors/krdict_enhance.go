package interceptors

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"launchpad.net/xmlpath"
)

func BackfillEnglishDefinition(item *pkg.Item, options map[string]string) error {
	if !strings.HasPrefix(item.Id, "krdict") {
		return nil
	}

	if item.Def.English != "" {
		return nil
	}

	// format: strings.Join([]string{"krdict", lang, entryId, lexicalUnit}, ":")
	entryId := strings.Split(item.Id, ":")[2]

	item.Def.English = fetchEnglishDefinition(entryId)

	return nil
}

func fetchEnglishDefinition(entryId string) string {
	if os.Getenv("KRDICT_API_KEY") == "" {
		panic("KRDICT_API_KEY not set.")
	}

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
	transWord := get(node, transPath+"/trans_word")
	transDfn := get(node, transPath+"/trans_dfn")

	return transWord + " := " + transDfn
}

func get(node *xmlpath.Node, xpath string) string {
	path := xmlpath.MustCompile(xpath)

	if value, ok := path.String(node); ok {
		return value
	}

	return ""
}

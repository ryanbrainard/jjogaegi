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

	if item.Def.English == "" || strings.TrimSpace(item.Def.English) == ":=" {
		item.Def.English = getEnglishDefinition(entry)
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

func get(node *xmlpath.Node, xpath string) string {
	path := xmlpath.MustCompile(xpath)

	if value, ok := path.String(node); ok {
		return value
	}

	return ""
}

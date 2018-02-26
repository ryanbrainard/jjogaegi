package interceptors

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"launchpad.net/xmlpath"
)

func KrDictLookup(item *pkg.Item, options map[string]string) error {
	if options[pkg.OPT_LOOKUP] != "true" {
		return nil
	}

	if strings.HasPrefix(item.ExternalID, "krdict") {
		return nil
	}

	q := item.Hangul
	if item.Hanja != "" {
		q = item.Hanja
	}

	results, err := search(q, options)
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

	var choiceIndex int
	switch len(choices) {
	case 0:
		lookupOut("No results found for %s. Skipping lookup.\n", item.Hangul)
		return nil
	case 1:
		choiceIndex = 0
	default:
		lookupOut("Multiple results found for %s:\n", item.Hangul)
		for i, choice := range choices {
			println("[lookup] ", i+1, ") ", pkg.XpathString(choice, "sense/translation/trans_word"))
		}
		if options[pkg.OPT_INTERACTIVE] == "true" {
			choiceIndex = promptMultipleChoice(item, choices)
		} else {
			lookupOut("Skipping lookup. Set %s option to choose.\n\n", pkg.OPT_INTERACTIVE)
			return nil
		}
	}

	item.ExternalID = pkg.KrDictID("kor", pkg.XpathString(choices[choiceIndex], "target_code"), "단어")
	return nil
}

func search(q string, options map[string]string) (*xmlpath.Node, error) {
	url := fmt.Sprintf(
		"%s/api/search?key=%s&type_search=search&part=word&q=%s&sort=dict&translated=y&trans_lang=1",
		options[pkg.OPT_KRDICT_API_URL],
		options[pkg.OPT_KRDICT_API_KEY],
		q,
	)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("search type=eng url=%q err=%q", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	return xmlpath.Parse(resp.Body)
}

func promptMultipleChoice(item *pkg.Item, choices []*xmlpath.Node) int {
	for {
		lookupOut("Enter number: ")
		answerString, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			lookupOut("%s\n", err)
			continue
		}

		answerNum, err := strconv.Atoi(strings.TrimSpace(answerString))
		if err != nil || answerNum < 1 || answerNum > len(choices) {
			lookupOut("Invalid number\n")
			continue
		}

		lookupOut("\n")
		return answerNum - 1
	}
}

func lookupOut(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "[lookup] "+format, a...)
}

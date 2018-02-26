package interceptors

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"launchpad.net/xmlpath"
)

func NewKrDictLookup(interactiveIn io.Reader, interactiveOut io.Writer) pkg.InterceptorFunc {
	return func(item *pkg.Item, options map[string]string) error {
		return krDictLookup(interactiveIn, interactiveOut, item, options)
	}
}

func krDictLookup(in io.Reader, out io.Writer, item *pkg.Item, options map[string]string) error {
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

	itemLabel := item.Hangul
	if item.Def.English != "" {
		itemLabel += " (" + item.Def.English + ")"
	}

	var choiceIndex int
	switch len(choices) {
	case 0:
		fmt.Fprintf(out, "No results found for %s. Skipping lookup.\n", itemLabel)
		return nil
	case 1:
		choiceIndex = 0
	default:
		fmt.Fprintf(out, "Multiple results found for %s:\n", itemLabel)
		for i, choice := range choices {
			fmt.Fprintf(out, " %d) %s\n", i+1, pkg.XpathString(choice, "sense/translation/trans_word"))
		}
		if options[pkg.OPT_INTERACTIVE] == "true" {
			choiceIndex = promptMultipleChoice(in, out, item, choices)
		} else {
			fmt.Fprintf(out, "Skipping lookup. Set %s option to choose.\n", pkg.OPT_INTERACTIVE)
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
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Non-200 searching KR DICT API")
	}

	return xmlpath.Parse(resp.Body)
}

func promptMultipleChoice(in io.Reader, out io.Writer, item *pkg.Item, choices []*xmlpath.Node) int {
	inBuf := bufio.NewReader(in)
	for {
		fmt.Fprintf(out, "Enter number: ")
		answerString, err := inBuf.ReadString('\n')
		if err != nil {
			fmt.Fprintf(out, "%s\n", err)
			continue
		}

		answerNum, err := strconv.Atoi(strings.TrimSpace(answerString))
		if err != nil || answerNum < 1 || answerNum > len(choices) {
			fmt.Fprintf(out, "Invalid number\n")
			continue
		}

		fmt.Fprintf(out, "\n")
		return answerNum - 1
	}
}

package parsers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"golang.org/x/net/html"
)

var callbackStartBytes = []byte("window.__jindo2_callback")
var callbackEndByte = byte('(')

func ParseNaverJSON(r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	buf := bufio.NewReader(r)
	header, err := buf.Peek(len(callbackStartBytes))
	if err != nil {
		return err
	}

	if string(header) == string(callbackStartBytes) {
		buf.ReadString(callbackEndByte)
	}

	dec := json.NewDecoder(buf)

	// read open bracket
	_, err = dec.Token()
	if err != nil {
		return err
	}

	for dec.More() {
		var page NaverPage
		// decode an array value (Message)
		err := dec.Decode(&page)
		if err != nil {
			return err
		}

		for _, item := range page.Items {
			hangulTerm, hanjaTerm := splitHangul(item.renderItem())

			examples := []pkg.Translation{}
			for _, means := range item.Means {
				for _, example := range means.Examples {
					examples = append(examples, pkg.Translation{
						English: stripHTML(example.English),
						Korean:  example.Korean,
					})
				}
			}

			items <- &pkg.Item{
				NoteID: strings.Join([]string{"naver", item.OriginEntryId}, ":"),
				Hangul: hangulTerm,
				Hanja:  hanjaTerm,
				Def: pkg.Translation{
					English: item.renderMeans()},
				Examples: examples,
			}
		}
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		return err
	}

	return nil
}

type NaverPage struct {
	Items []NaverItem `json:"items"`
}

type NaverItem struct {
	OriginEntryId string      `json:"originEntryID"`
	EntryName     string      `json:"entryName"`
	Means         []NaverMean `json:"means"`
}

func (i NaverItem) renderItem() string {
	return stripHTML(i.EntryName)
}

func (i NaverItem) renderMeans() string {
	renderedMeans := []string{}
	for _, m := range i.Means {
		rm := ""
		if len(i.Means) > 1 {
			rm = fmt.Sprintf("%d. ", m.Seq)
		}
		rm += m.render()
		renderedMeans = append(renderedMeans, rm)
	}
	return strings.Join(renderedMeans, "  ")
}

type NaverMean struct {
	Seq      int            `json:"seq"`
	Mean     string         `json:"mean"`
	Examples []NaverExample `json:"examples"`
}

type NaverExample struct {
	English string `json:"example"`
	Korean  string `json:"translated"`
}

func (m NaverMean) render() string {
	return stripHTML(m.Mean)
}

func stripHTML(in string) string {
	z := html.NewTokenizer(strings.NewReader(in))
	out := ""
	for {
		switch z.Next() {
		case html.TextToken:
			out += string(z.Text())
		case html.ErrorToken:
			return out
		}
	}
	return out
}

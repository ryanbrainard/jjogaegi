package parsers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"golang.org/x/net/html"
	"io"
	"strings"
)

var callbackStartBytes = []byte("window.__jindo2_callback")
var callbackEndByte = byte('(')

func ParseNaverJSON(r io.Reader, items chan<- *pkg.Item) {
	buf := bufio.NewReader(r)
	header, err := buf.Peek(len(callbackStartBytes))
	if err != nil {
		panic(err)
	}

	if string(header) == string(callbackStartBytes) {
		buf.ReadString(callbackEndByte)
	}

	dec := json.NewDecoder(buf)

	// read open bracket
	_, err = dec.Token()
	if err != nil {
		panic(err)
	}

	for dec.More() {
		var page NaverPage
		// decode an array value (Message)
		err := dec.Decode(&page)
		if err != nil {
			panic(err)
		}

		for _, item := range page.Items {
			hangeulTerm, hanjaTerm := splitHangeul(item.EntryName)
			items <- &pkg.Item{
				Term:    hangeulTerm,
				SubTerm: hanjaTerm,
				Def:     item.renderMeans(),
			}
		}
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		panic(err)
	}

	close(items)
}

type NaverPage struct {
	Items []NaverItem `json:"items"`
}

type NaverItem struct {
	EntryName string      `json:"entryName"`
	Means     []NaverMean `json:"means"`
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
	Seq  int    `json:"seq"`
	Mean string `json:"mean"`
}

func (m NaverMean) render() string {
	z := html.NewTokenizer(strings.NewReader(m.Mean))
	s := ""
	for {
		switch z.Next() {
		case html.TextToken:
			s += string(z.Text())
		case html.ErrorToken:
			return s
		}
	}
	return s
}

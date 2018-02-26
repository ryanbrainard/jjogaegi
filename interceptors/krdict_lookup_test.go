package interceptors

import (
	"bytes"
	"os"
	"strconv"
	"testing"

	"github.com/ryanbrainard/jjogaegi/pkg"
	mocks "github.com/ryanbrainard/jjogaegi/testing"
	"github.com/stretchr/testify/assert"
)

type krDictLookupTestCase struct {
	name             string
	item             *pkg.Item
	expectedItem     *pkg.Item
	expectedOut      string
	lookup           bool
	interactive      bool
	interactiveInput string
}

func TestKrDictLookup(t *testing.T) {
	ts := mocks.NewKrdictMockServer()
	defer ts.Close()

	cases := []krDictLookupTestCase{
		{
			name:         "multiple results/lookup/non-interactive",
			lookup:       true,
			interactive:  false,
			item:         &pkg.Item{Hangul: "안녕"},
			expectedItem: &pkg.Item{Hangul: "안녕"},
			expectedOut:  "Multiple results found for 안녕:\n 1) hello; hi; good-bye; bye\n 2) peace; good health\nSkipping lookup. Set interactive option to choose.\n\n",
		},
		{
			name:             "multiple results/lookup/interactive",
			lookup:           true,
			interactive:      true,
			item:             &pkg.Item{Hangul: "안녕"},
			interactiveInput: "2\n",
			expectedItem:     &pkg.Item{Hangul: "안녕", ExternalID: "krdict:kor:17296:단어"},
			expectedOut:      "Multiple results found for 안녕:\n 1) hello; hi; good-bye; bye\n 2) peace; good health\nEnter number: \n",
		},
		{
			name:             "multiple results/lookup/interactive/bad response",
			lookup:           true,
			interactive:      true,
			item:             &pkg.Item{Hangul: "안녕"},
			interactiveInput: "X\n2\n",
			expectedItem:     &pkg.Item{Hangul: "안녕", ExternalID: "krdict:kor:17296:단어"},
			expectedOut:      "Multiple results found for 안녕:\n 1) hello; hi; good-bye; bye\n 2) peace; good health\nEnter number: Invalid number\nEnter number: \n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(rt *testing.T) {
			options := map[string]string{
				pkg.OPT_KRDICT_API_URL: ts.URL,
				pkg.OPT_KRDICT_API_KEY: os.Getenv("KRDICT_API_KEY"),
				pkg.OPT_LOOKUP:         strconv.FormatBool(c.lookup),
				pkg.OPT_INTERACTIVE:    strconv.FormatBool(c.interactive),
			}

			in := bytes.NewBufferString(c.interactiveInput)
			out := &bytes.Buffer{}
			actual := c.item

			err := NewKrDictLookup(in, out)(actual, options)

			assert.NoError(t, err)
			assert.Equal(t, c.expectedItem, actual)
			assert.Equal(t, out.String(), c.expectedOut)
		})
	}
}

package parsers

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/ryanbrainard/jjogaegi/pkg"
	mocks "github.com/ryanbrainard/jjogaegi/testing"
	"github.com/stretchr/testify/assert"
)

type promptTestCase struct {
	name             string
	interactiveInput string
	expectedItems    []*pkg.Item
	expectedOut      string
}

func TestPrompt(t *testing.T) {
	ts := mocks.NewKrdictMockServer()
	defer ts.Close()

	cases := []promptTestCase{
		{
			name:             "korean only",
			interactiveInput: "안경\n",
			expectedItems:    []*pkg.Item{{Hangul: "안경", ExternalID: "krdict:kor:31484:단어"}},
			expectedOut:      "Enter a Korean word on each line: (press Ctrl+D to quit)\n>>> 안경 -> glasses; spectacles\n\n>>> ",
		},
		{
			name:             "korean and english",
			interactiveInput: "안경 specs\n",
			expectedItems:    []*pkg.Item{{Hangul: "안경", ExternalID: "krdict:kor:31484:단어"}},
			expectedOut:      "Enter a Korean word on each line: (press Ctrl+D to quit)\n>>> 안경 -> glasses; spectacles\n\n>>> ",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(rt *testing.T) {
			options := map[string]string{
				pkg.OPT_MEDIADIR:       "DUMMY",
				pkg.OPT_KRDICT_API_URL: ts.URL,
				pkg.OPT_KRDICT_API_KEY: os.Getenv("KRDICT_API_KEY"),
			}

			if options[pkg.OPT_KRDICT_API_KEY] == "" {
				options[pkg.OPT_KRDICT_API_KEY] = "DUMMY"
			}

			in := bytes.NewBufferString(c.interactiveInput)
			out := bytes.NewBufferString("")
			items := make(chan *pkg.Item, 100)

			err := NewInteractivePrompt(out)(context.Background(), in, items, options)

			assert.NoError(t, err)
			assert.Equal(t, c.expectedOut, out.String())
			for _, expectedItem := range c.expectedItems {
				assert.Equal(t, expectedItem, <-items)
			}
		})
	}
}

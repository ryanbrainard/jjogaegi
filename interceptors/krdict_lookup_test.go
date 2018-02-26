package interceptors

import (
	"bytes"
	"os"
	"testing"

	"github.com/ryanbrainard/jjogaegi/pkg"
	mocks "github.com/ryanbrainard/jjogaegi/testing"
	"github.com/stretchr/testify/assert"
)

func TestKrDictLookup(t *testing.T) {
	// TODO: consider making this also hit real server
	ts := mocks.NewKrdictMockServer()
	defer ts.Close()
	options := map[string]string{
		pkg.OPT_KRDICT_API_URL: ts.URL,
		pkg.OPT_KRDICT_API_KEY: os.Getenv("KRDICT_API_KEY"),
		pkg.OPT_LOOKUP:         "true",
	}

	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	actualItem := &pkg.Item{Hangul: "안녕"}
	expectedItem := &pkg.Item{Hangul: "안녕"}

	err := NewKrDictLookup(in, out)(actualItem, options)
	assert.NoError(t, err)
	assert.Equal(t, expectedItem, actualItem)
	assert.Equal(t, out.String(), "Multiple results found for 안녕:\n 1) hello; hi; good-bye; bye\n 2) peace; good health\nSkipping lookup. Set interactive option to choose.\n\n")
}

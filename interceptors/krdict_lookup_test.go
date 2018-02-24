package interceptors

import (
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
		pkg.OPT_LOOKUP:         "true",
	}

	actualItem := &pkg.Item{Hangul: "안녕"}
	expectedItem := &pkg.Item{Hangul: "안녕"}

	err := KrDictLookup(actualItem, options)
	assert.NoError(t, err)
	assert.Equal(t, expectedItem, actualItem)
}

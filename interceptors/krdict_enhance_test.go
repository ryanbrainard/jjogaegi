package interceptors

import (
	"os"
	"testing"

	"launchpad.net/xmlpath"

	"github.com/ryanbrainard/jjogaegi/pkg"
	mocks "github.com/ryanbrainard/jjogaegi/testing"
	"github.com/stretchr/testify/assert"
)

func TestKrDictEnhance(t *testing.T) {
	// TODO: consider making this also hit real server
	ts := mocks.NewKrdictMockServer()
	defer ts.Close()
	options := map[string]string {
		pkg.OPT_KRDICT_API_URL: ts.URL,
	}

	for i, expectedItem := range xmlTestItems {
		t.Run(expectedItem.Hangul, func(tr *testing.T) {
			// TODO: remove skips
			if i != 0 {
				tr.Skip("TODO: add mock")
			}

			actualItem := &pkg.Item{
				ExternalID: expectedItem.ExternalID,
			}

			err := KrDictEnhance(actualItem, options)
			assert.NoError(t, err)
			assert.Equal(tr, expectedItem, actualItem)
		})
	}
}

// TODO: stuff into mocks and delete this test
func TestGetWordGrade(t *testing.T) {
	in, err := os.Open("../testing/fixtures/kr_dict_en_15392-mod.xml")
	assert.Nil(t, err)
	node, err := xmlpath.Parse(in)
	assert.Nil(t, err)

	assert.Equal(t, "고급", getWordGrade(node))
}

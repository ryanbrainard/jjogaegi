package interceptors

import (
	"os"
	"testing"

	"launchpad.net/xmlpath"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

func TestKrDictEnhance(t *testing.T) {
	if os.Getenv("KRDICT_API_KEY") == "" {
		t.Skip("KRDICT_API_KEY not set")
	}

	for _, expectedItem := range xmlTestItems {
		t.Run(expectedItem.Hangul, func(tr *testing.T) {
			actualItem := &pkg.Item{
				ExternalID: expectedItem.ExternalID,
			}

			err := KrDictEnhance(actualItem, map[string]string{})
			assert.NoError(t, err)
			assert.Equal(tr, expectedItem, actualItem)
		})
	}
}
func TestGetEnglishDefinition(t *testing.T) {
	in, err := os.Open("../parsers/fixtures/kr_dict_en_15392.xml")
	assert.Nil(t, err)
	node, err := xmlpath.Parse(in)
	assert.Nil(t, err)

	assert.Equal(t, "join; sign up := To join a group or sign up for goods and services.", getEnglishDefinition(node))
}
func TestGetWordGrade(t *testing.T) {
	in, err := os.Open("../parsers/fixtures/kr_dict_en_15392-mod.xml")
	assert.Nil(t, err)
	node, err := xmlpath.Parse(in)
	assert.Nil(t, err)

	assert.Equal(t, "고급", getWordGrade(node))
}

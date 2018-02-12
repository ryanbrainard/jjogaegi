package interceptors

import (
	"os"
	"testing"

	"launchpad.net/xmlpath"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

func TestKrDictEnhance(t *testing.T) {
	item := &pkg.Item{
		ExternalID: "krdict:kor:15392:단어",
	}

	err := KrDictEnhance(item, map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, "join; sign up := To join a group or sign up for goods and services.", item.Def.English)
}
func TestExtractEnglishDefinition(t *testing.T) {
	in, err := os.Open("../parsers/fixtures/kr_dict_en_15392.xml")
	assert.Nil(t, err)
	node, err := xmlpath.Parse(in)
	assert.Nil(t, err)

	assert.Equal(t, "join; sign up := To join a group or sign up for goods and services.", getEnglishDefinition(node))
}

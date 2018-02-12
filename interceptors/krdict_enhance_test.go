package interceptors

import (
	"os"
	"testing"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

func TestBackfilEnglishDefinition(t *testing.T) {
	item := &pkg.Item{
		Id: "krdict:kor:15392:단어",
	}

	err := BackfillEnglishDefinition(item, map[string]string{})
	assert.NoError(t, err)
	assert.Equal(t, "join; sign up := To join a group or sign up for goods and services.", item.Def.English)
}
func TestExtractEnglishDefinition(t *testing.T) {
	in, err := os.Open("../parsers/fixtures/kr_dict_en_15392.xml")
	assert.Nil(t, err)

	assert.Equal(t, "join; sign up := To join a group or sign up for goods and services.", extractEnglishDefinition(in))
}

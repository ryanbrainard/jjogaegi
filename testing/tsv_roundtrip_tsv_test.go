package testing

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"

	"ryanbrainard.com/jjogaegi/formatters"
	"ryanbrainard.com/jjogaegi/parsers"
	"ryanbrainard.com/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

func TestRoundtripTSV(t *testing.T) {
	file, err := os.Open("../testing/fixtures/sample-1.tsv")
	assert.Nil(t, err)

	inBytes, err := ioutil.ReadAll(file)
	assert.Nil(t, err)

	in := bytes.NewBuffer(inBytes)

	items := make(chan *pkg.Item, 100)
	err = parsers.ParseTSV(context.Background(), in, items, map[string]string{})
	assert.Nil(t, err)

	close(items)

	out := &bytes.Buffer{}
	err = formatters.FormatTSV(context.Background(), items, out, map[string]string{})
	assert.Nil(t, err)

	assert.Equal(t, string(inBytes), out.String())
}

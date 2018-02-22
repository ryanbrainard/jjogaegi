package interceptors

import (
	"testing"

	"github.com/ryanbrainard/jjogaegi/pkg"
	mocks "github.com/ryanbrainard/jjogaegi/testing"

	"github.com/stretchr/testify/assert"

	"strings"
)

func TestMediaFormatting(t *testing.T) {
	// TODO: consider making this also hit real server
	ts := mocks.NewKrdictMockServer()
	defer ts.Close()

	item := xmlTestItems[2]
	item.ImageTag = strings.Replace(item.ImageTag, "http://dicmedia.korean.go.kr:8899", ts.URL, 1)

	err := MediaFormatting(item, map[string]string{pkg.OPT_MEDIADIR: "/tmp"})
	assert.NoError(t, err)
	assert.Equal(t, "", item.AudioTag) // TODO: fill in
	assert.Regexp(t, "<img src=\"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}.jpg\">", item.ImageTag)
}

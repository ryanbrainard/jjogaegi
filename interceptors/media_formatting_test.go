package interceptors

import (
	"testing"

	"ryanbrainard.com/jjogaegi/pkg"
	mocks "ryanbrainard.com/jjogaegi/testing"

	"github.com/stretchr/testify/assert"

	"strings"
)

func TestMediaFormatting(t *testing.T) {
	ts := mocks.NewKrdictMockServer()
	defer ts.Close()

	const baseUrl = "http://dicmedia.korean.go.kr:8899"
	item := xmlTestItems[2]
	item.ImageTag = strings.Replace(item.ImageTag, baseUrl, ts.URL, 1)
	item.ImageTag += "&_baseUrl=" + baseUrl

	err := MediaFormatting(item, map[string]string{pkg.OPT_MEDIADIR: "/tmp", pkg.OPT_DEBUG: "true"})
	assert.NoError(t, err)
	assert.Equal(t, "", item.AudioTag) // TODO: fill in
	assert.Regexp(t, "<img src=\"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}.jpg\">", item.ImageTag)
}

func TestMediaFormattingWithoutMediaDir(t *testing.T) {
	item := &pkg.Item{ImageTag: "http://example.com/image.jpg"}

	err := MediaFormatting(item, map[string]string{})
	assert.EqualError(t, err, "cannot download media (http://example.com/image.jpg) unless media dir is set")
}

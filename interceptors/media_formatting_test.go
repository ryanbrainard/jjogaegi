package interceptors

import (
	"testing"

	"github.com/ryanbrainard/jjogaegi/pkg"

	"github.com/stretchr/testify/assert"
)

func TestMediaFormatting(t *testing.T) {
	item := xmlTestItems[2]
	err := MediaFormatting(item, map[string]string{pkg.OPT_MEDIADIR: "/tmp"})
	assert.NoError(t, err)
	assert.Equal(t, "", item.AudioTag) // TODO: fill in
	assert.Regexp(t, "<img src=\"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}.jpg\">", item.ImageTag)
}

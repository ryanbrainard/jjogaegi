package formatters

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatJSON(t *testing.T) {
	items, out := setupTestFormat()
	err := FormatJSON(context.Background(), items, out, map[string]string{})
	assert.Nil(t, err)
	assert.Equal(t, "[\n{\n  \"NoteID\": \"\",\n  \"ExternalID\": \"\",\n  \"Hangul\": \"처리\",\n  \"Hanja\": \"處理\",\n  \"Def\": {\n    \"Korean\": \"\",\n    \"English\": \"handling\"\n  },\n  \"Pronunciation\": \"\",\n  \"AudioTag\": \"\",\n  \"ImageTag\": \"\",\n  \"Grade\": \"\",\n  \"Antonym\": \"\",\n  \"Examples\": [\n    {\n      \"Korean\": \"k\",\n      \"English\": \"e\"\n    }\n  ]\n}\n]\n", out.String())
}

package jjogaegi

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	in := strings.NewReader(" • 안녕 라이언 Hello, Ryan")
	out := &bytes.Buffer{}
	Parse(in, out)
	assert.Equal(t, "안녕 라이언\tHello, Ryan\n", out.String())
}

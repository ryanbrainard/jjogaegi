package pkg

import (
	"context"
	"io"
)

type Item struct {
	NoteID        string
	ExternalID    string
	Hangul        string
	Hanja         string
	Def           Translation
	Pronunciation string
	AudioTag      string
	ImageTag      string
	Grade         string
	Antonym       string
	Examples      []Translation
}

type Translation struct {
	Korean  string
	English string
}

type ParseFunc func(ctx context.Context, reader io.Reader, items chan<- *Item, options map[string]string) error

type FormatFunc func(ctx context.Context, items <-chan *Item, writer io.Writer, options map[string]string) error

type InterceptorFunc func(item *Item, options map[string]string) error

const OPT_HANJA = "hanja"
const OPT_HANJA_NONE = "none"
const OPT_HANJA_PARENTHESIS = "parens"
const OPT_HEADER = "header"
const OPT_MEDIADIR = "mediadir"
const OPT_PARALLEL = "parallel"
const OPT_LOOKUP = "lookup"
const OPT_INTERACTIVE = "interactive"
const OPT_KRDICT_API_URL = "krdict_api_url"
const OPT_KRDICT_API_KEY = "krdict_api_key"

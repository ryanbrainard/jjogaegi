package pkg

import (
	"io"
)

type Item struct {
	NoteId        string
	Id            string
	Hangul        string
	Hanja         string
	Pronunciation string
	AudioURL      string // TODO: drop
	AudioTag      string // TODO: drop
	ImageURL      string
	ImageTag      string
	Def           Translation
	Antonym       string
	Examples      []Translation
}

type Translation struct {
	Korean  string
	English string
}

type ParseFunc func(reader io.Reader, items chan<- *Item, options map[string]string) error

type FormatFunc func(items <-chan *Item, writer io.Writer, options map[string]string) error

const OPT_HANJA = "hanja"
const OPT_HANJA_NONE = "none"
const OPT_HANJA_PARENTHESIS = "parens"
const OPT_HEADER = "header"
const OPT_MEDIADIR = "mediadir"

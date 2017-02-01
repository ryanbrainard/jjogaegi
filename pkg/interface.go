package pkg

import "io"

type Item struct {
	Hangul string
	Hanja  string
	Def    string
}

type ParseFunc func(reader io.Reader, items chan<- *Item, options map[string]string)

type FormatFunc func(items <-chan *Item, writer io.Writer, options map[string]string)

const OPT_HANJA = "hanja"
const OPT_HANJA_NONE = "none"
const OPT_HANJA_PARENTHESIS = "parens"
const OPT_HANJA_SEPARATE = "sep"
const OPT_HEADER = "header"

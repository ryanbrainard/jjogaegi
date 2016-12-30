package pkg

import "io"

type Item struct {
	Hangul string
	Hanja  string
	Def    string
}

type ParseFunc func(reader io.Reader, items chan<- *Item)

type FormatFunc func(items <-chan *Item, writer io.Writer)

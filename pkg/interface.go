package pkg

import "io"

type Item struct {
	Term    string
	SubTerm string
	Def     string
}

type ParseFunc func(reader io.Reader, items chan<- *Item)

type FormatFunc func(items <-chan *Item, writer io.Writer)

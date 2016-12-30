package pkg

import (
	"io"
)

func Run(in io.Reader, parse ParseFunc, format FormatFunc, out io.Writer) {
	items := make(chan *Item)
	go parse(in, items)
	format(items, out)
}

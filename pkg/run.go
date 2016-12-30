package pkg

import (
	"io"
)

func Run(in io.Reader, out io.Writer, parse ParseFunc, format FormatFunc, options map[string]string) {
	items := make(chan *Item)
	go parse(in, items, options)
	format(items, out, options)
}

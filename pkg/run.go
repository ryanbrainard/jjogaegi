package pkg

import (
	"fmt"
	"io"
)

func Run(in io.Reader, out io.Writer, parse ParseFunc, format FormatFunc, options map[string]string) error {
	if parse == nil {
		return fmt.Errorf("Missing or invalid parser specified")
	}

	if format == nil {
		return fmt.Errorf("Missing or invalid formatter specified")
	}

	items := make(chan *Item)
	var err error

	go func() {
		err = parse(in, items, options)
		close(items)
	}()

	if err := format(items, out, options); err != nil {
		return err
	}

	return err
}

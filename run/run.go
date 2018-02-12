package run

import (
	"fmt"
	"io"

	"github.com/ryanbrainard/jjogaegi/interceptors"
	"github.com/ryanbrainard/jjogaegi/pkg"
)

func Run(in io.Reader, out io.Writer, parse pkg.ParseFunc, format pkg.FormatFunc, options map[string]string) error {
	if parse == nil {
		return fmt.Errorf("Missing or invalid parser specified")
	}

	if format == nil {
		return fmt.Errorf("Missing or invalid formatter specified")
	}

	var err error // TODO: handle this with a channel??

	parsed := make(chan *pkg.Item)
	go func() {
		err = parse(in, parsed, options)
		close(parsed)
	}()

	interceptors := []pkg.InterceptorFunc{
		interceptors.BackfillEnglishDefinition,
	}
	intercepted := make(chan *pkg.Item)
	go func() {
		for item := range parsed {
			for _, interceptor := range interceptors {
				if err != nil {
					break
				}
				err = interceptor(item, options)
			}
			intercepted <- item
		}
		close(intercepted)
	}()

	// TODO: error handle

	if err := format(intercepted, out, options); err != nil {
		return err
	}

	return err
}

package parsers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/ryanbrainard/jjogaegi/interceptors"
	"github.com/ryanbrainard/jjogaegi/pkg"
)

func NewInteractivePrompt(interactiveOut io.Writer) pkg.ParseFunc {
	return func(ctx context.Context, reader io.Reader, items chan<- *pkg.Item, options map[string]string) error {
		return interactivePrompt(interactiveOut, ctx, reader, items, options)
	}
}

func interactivePrompt(interactiveOut io.Writer, ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	options[pkg.OPT_LOOKUP] = strconv.FormatBool(true)
	options[pkg.OPT_INTERACTIVE] = strconv.FormatBool(true)

	fmt.Fprintf(interactiveOut, "Enter a Korean word on each line: (press Ctrl+D to quit)\n")
	prompt := ">>> "
	fmt.Fprintf(interactiveOut, prompt)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := sanitize(scanner.Text())

		if line == "" {
			fmt.Fprintf(interactiveOut, "\n"+prompt)
			continue
		}

		item := parseLineItem(line)

		// pre-run interceptor to not muck up stdin processing
		err := interceptors.NewKrDictLookup(r, interactiveOut)(item, options)
		if err != nil {
			return err
		}

		items <- item

		// give time to allow interceptors to error and cancel before re-prompting
		time.Sleep(500 * time.Millisecond) // TODO: can we do this without sleeping
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fmt.Fprintf(interactiveOut, "\n"+prompt)
		}
	}

	return nil
}

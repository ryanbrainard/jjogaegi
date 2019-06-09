package parsers

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"go.ryanbrainard.com/jjogaegi/interceptors"
	"go.ryanbrainard.com/jjogaegi/pkg"
)

func NewInteractivePrompt(interactiveOut io.Writer) pkg.ParseFunc {
	return func(ctx context.Context, reader io.Reader, items chan<- *pkg.Item, options map[string]string) error {
		return interactivePrompt(interactiveOut, ctx, reader, items, options)
	}
}

func interactivePrompt(interactiveOut io.Writer, ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	if options[pkg.OPT_KRDICT_API_KEY] == "" || options[pkg.OPT_MEDIADIR] == "" {
		return fmt.Errorf("KRDICT_API_KEY and MEDIA_DIR must be set in environment. Run `jjogaegi -help` for details.")
	}

	options[pkg.OPT_LOOKUP] = strconv.FormatBool(true)
	options[pkg.OPT_INTERACTIVE] = strconv.FormatBool(true)

	fmt.Fprintf(interactiveOut, "Enter a Korean word on each line: (press Ctrl+D to quit)\n")
	prompt := ">>> "
	fmt.Fprintf(interactiveOut, prompt)

	lookup := interceptors.NewKrDictLookup(r, interactiveOut)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := sanitize(scanner.Text())
		hangul, _ := splitHangul(line)

		if hangul == "" {
			fmt.Fprintf(interactiveOut, "<invalid input>\n%s", prompt)
			continue
		}

		item := &pkg.Item{
			Hangul: sanitize(hangul),
		}

		// pre-run interceptor to not muck up stdin processing
		err := lookup(item, options)
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
			fmt.Fprint(interactiveOut, prompt)
		}
	}

	return nil
}

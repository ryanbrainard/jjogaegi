package parsers

import (
	"bufio"
	"io"
	"os"
	"time"

	"github.com/ryanbrainard/jjogaegi/interceptors"
	"github.com/ryanbrainard/jjogaegi/pkg"
)

func InteractivePrompt(r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	options[pkg.OPT_LOOKUP] = "true"
	options[pkg.OPT_INTERACTIVE] = "true"

	println("Enter a Korean word on each line: (press Ctrl+D to quit)")
	prompt := ">>> "
	print(prompt)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := sanitize(scanner.Text())

		if line == "" {
			print("\n" + prompt)
			continue
		}

		item := &pkg.Item{
			Hangul: sanitize(line),
		}

		// pre-run interceptor to not muck up stdin processing
		err := interceptors.NewKrDictLookup(os.Stdin, os.Stderr)(item, options)
		if err != nil {
			return err
		}

		items <- item

		time.Sleep(200 * time.Millisecond) // TODO: sync to remove this?
		print("\n" + prompt)
	}

	return nil
}

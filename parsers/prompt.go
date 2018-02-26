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
	println("Enter a Korean word on each line:")
	print(">>> ")

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		item := &pkg.Item{
			Hangul: sanitize(line),
		}

		// pre-run interceptor to not muck up stdin processing
		options[pkg.OPT_LOOKUP] = "true"
		err := interceptors.NewKrDictLookup(os.Stdin, os.Stderr)(item, options)
		if err != nil {
			return err
		}

		items <- item

		time.Sleep(1 * time.Second)
		println()
		print(">>> ")
	}

	return nil
}

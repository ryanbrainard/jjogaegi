package run

import (
	"fmt"
	"io"
	"runtime"
	"sync"

	"github.com/ryanbrainard/jjogaegi/interceptors"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"golang.org/x/sync/errgroup"
)

func Run(in io.Reader, out io.Writer, parse pkg.ParseFunc, format pkg.FormatFunc, options map[string]string) error {
	if parse == nil {
		return fmt.Errorf("Missing or invalid parser specified")
	}

	if format == nil {
		return fmt.Errorf("Missing or invalid formatter specified")
	}

	interceptors := []pkg.InterceptorFunc{
		interceptors.GenerateNoteId,
		interceptors.KrDictEnhance,
	}

	parallelism := 1
	if options[pkg.OPT_PARALLEL] == "true" {
		parallelism = runtime.NumCPU()
	}

	var g errgroup.Group

	parsed := make(chan *pkg.Item)
	g.Go(func() error {
		err := parse(in, parsed, options)
		close(parsed)
		return err
	})

	var iwg sync.WaitGroup
	iwg.Add(parallelism)

	intercepted := make(chan *pkg.Item)
	for p := 0; p < parallelism; p++ {
		g.Go(func() error {
			for item := range parsed {
				for _, interceptor := range interceptors {
					if err := interceptor(item, options); err != nil {
						return err
					}
				}
				intercepted <- item
			}
			iwg.Done()
			return nil
		})
	}

	go func() {
		iwg.Wait()
		close(intercepted)
	}()

	if err := format(intercepted, out, options); err != nil {
		return err
	}

	return g.Wait()
}

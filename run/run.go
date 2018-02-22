package run

import (
	"fmt"
	"io"
	"runtime"
	"sync"

	"github.com/ryanbrainard/jjogaegi/interceptors"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"golang.org/x/sync/errgroup"
	"os"
)

func Run(in io.Reader, out io.Writer, parse pkg.ParseFunc, format pkg.FormatFunc, options map[string]string) error {
	if parse == nil {
		return fmt.Errorf("Missing or invalid parser specified")
	}

	if format == nil {
		return fmt.Errorf("Missing or invalid formatter specified")
	}

	parallelism := 1
	if options[pkg.OPT_PARALLEL] == "true" {
		parallelism = runtime.NumCPU()
	}

	setEnvOpt(options, pkg.OPT_KRDICT_API_KEY, "KRDICT_API_KEY", "")
	setEnvOpt(options, pkg.OPT_KRDICT_API_URL, "KRDICT_API_URL", "https://krdict.korean.go.kr")
	setEnvOpt(options, pkg.OPT_MEDIADIR, "ANKI_MEDIA_DIR", "")

	var g errgroup.Group

	parsed := make(chan *pkg.Item)
	g.Go(func() error {
		err := parse(in, parsed, options)
		close(parsed)
		return err
	})

	interceptors := []pkg.InterceptorFunc{
		interceptors.GenerateNoteId,
		interceptors.KrDictLookup,
		interceptors.KrDictEnhance,
		interceptors.MediaFormatting,
	}
	intercepted := make(chan *pkg.Item)
	var iwg sync.WaitGroup
	iwg.Add(parallelism)
	for p := 0; p < parallelism; p++ {
		g.Go(func() error {
			defer iwg.Done()
			for item := range parsed {
				for _, interceptor := range interceptors {
					if err := interceptor(item, options); err != nil {
						return err
					}
				}
				intercepted <- item
			}
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

func setEnvOpt(options map[string]string, optKey, envKey, orDefault string) {
	if options[optKey] == "" {
		if envValue := os.Getenv(envKey); envValue != "" {
			options[optKey] = envValue
		} else {
			options[optKey] = orDefault
		}
	}
}

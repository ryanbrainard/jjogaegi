package run

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"sync"

	"os"

	"go.ryanbrainard.com/jjogaegi/interceptors"
	"go.ryanbrainard.com/jjogaegi/pkg"
	"golang.org/x/sync/errgroup"
)

func Run(in io.Reader, out io.Writer, parse pkg.ParseFunc, format pkg.FormatFunc, options map[string]string) error {
	if options == nil {
		options = map[string]string{}
	}

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

	setEnvOptOrDefault(options, pkg.OPT_KRDICT_API_URL, "KRDICT_API_URL", "https://krdict.korean.go.kr")
	setEnvOptOrDefault(options, pkg.OPT_KRDICT_API_KEY, "KRDICT_API_KEY", "")
	setEnvOptOrDefault(options, pkg.OPT_MEDIADIR, "MEDIA_DIR", "")

	pkg.Debug(options, "fn=Run at=group.create")
	g, ctx := errgroup.WithContext(context.Background())

	parsed := make(chan *pkg.Item)
	g.Go(func() error {
		err := parse(ctx, in, parsed, options)
		close(parsed)
		pkg.Debug(options, "fn=Run at=parse.closed err=%v", err)
		return err
	})

	interceptors := []pkg.InterceptorFunc{
		interceptors.GenerateNoteId,
		interceptors.NewKrDictLookup(os.Stdin, os.Stderr),
		interceptors.KrDictEnhance,
		interceptors.KrDictEnhanceHTML,
		interceptors.MediaFormatting,
	}
	intercepted := make(chan *pkg.Item)
	var iwg sync.WaitGroup
	iwg.Add(parallelism)
	for p := 0; p < parallelism; p++ {
		g.Go(func() error {
			defer iwg.Done()
			for item := range parsed {
				for i, interceptor := range interceptors {
					select {
					case <-ctx.Done():
						pkg.Debug(options, "fn=Run at=ctx.done err=%v", ctx.Err())
						return ctx.Err()
					default:
					}

					pkg.Debug(options, "fn=Run at=interceptor[%d].start hangul=%s", i, item.Hangul)
					if err := interceptor(item, options); err != nil {
						pkg.Debug(options, "fn=Run at=interceptor[%d].error hangul=%s error=%s", i, item.Hangul, err.Error())
						return err
					}
					pkg.Debug(options, "fn=Run at=interceptor[%d].done hangul=%s", i, item.Hangul)

				}
				intercepted <- item
			}
			pkg.Debug(options, "fn=Run at=interceptor.done err=%v", ctx.Err())
			return nil
		})
	}

	go func() {
		pkg.Debug(options, "fn=Run at=intercepted.waiting")
		iwg.Wait()
		close(intercepted)
		pkg.Debug(options, "fn=Run at=intercepted.close")
	}()

	if h, ok := options[pkg.OPT_HEADER]; ok && h != "" {
		out.Write([]byte(h + "\n"))
	}

	pkg.Debug(options, "fn=Run at=format.start")
	if err := format(ctx, intercepted, out, options); err != nil {
		return err
	}
	pkg.Debug(options, "fn=Run at=format.done")

	pkg.Debug(options, "fn=Run at=group.wait")
	return g.Wait()
}

func setEnvOptOrDefault(options map[string]string, optKey, envKey, orDefault string) {
	if options[optKey] == "" {
		if envValue := os.Getenv(envKey); envValue != "" {
			options[optKey] = envValue
		} else {
			options[optKey] = orDefault
		}
	}
}

func setEnvOptRequired(options map[string]string, optKey, envKey string) error {
	if options[optKey] == "" && os.Getenv(envKey) == "" {
		return fmt.Errorf("%s env must be set", envKey)

	}
	setEnvOptOrDefault(options, optKey, envKey, "")
	return nil
}

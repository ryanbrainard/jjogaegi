package run

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"ryanbrainard.com/jjogaegi/pkg"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	parser := func(ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
		items <- &pkg.Item{Hangul: "시험"}
		return nil
	}
	formatter := func(ctx context.Context, items <-chan *pkg.Item, w io.Writer, options map[string]string) error {
		for item := range items {
			w.Write([]byte(item.Hangul))
		}
		return nil
	}
	options := map[string]string{
		pkg.OPT_HEADER: "HDR",
	}

	err := Run(in, out, parser, formatter, options)
	assert.Nil(t, err)
	assert.Equal(t, "HDR\n시험", out.String())
}

func TestRun_ParserError(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	parser := func(ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
		return fmt.Errorf("boom: parser")
	}
	formatter := func(ctx context.Context, items <-chan *pkg.Item, w io.Writer, options map[string]string) error {
		for item := range items {
			w.Write([]byte(item.Hangul))
		}
		return nil
	}
	options := map[string]string{}

	err := Run(in, out, parser, formatter, options)
	assert.NotNil(t, err)
	assert.Equal(t, "", out.String())
}

func TestRun_InterceptorError(t *testing.T) {
	os.Setenv("MEDIA_DIR", "")
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	parser := func(ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
		inputItems := []*pkg.Item{
			{Hangul: "시험1", ImageTag: "http://example.com/image1.jpg"},
			{Hangul: "시험2", ImageTag: "http://example.com/image2.jpg"},
		}

		for _, ii := range inputItems {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				items <- ii
			}
			time.Sleep(100 * time.Millisecond) // TODO: yuck!
		}

		return nil
	}
	formatter := func(ctx context.Context, items <-chan *pkg.Item, w io.Writer, options map[string]string) error {
		for item := range items {
			w.Write([]byte(item.Hangul))
		}
		return nil
	}
	options := map[string]string{}

	err := Run(in, out, parser, formatter, options)
	assert.EqualError(t, err, "cannot download media (http://example.com/image1.jpg) unless media dir is set")
}

func TestRun_FormatterError(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	parser := func(ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
		items <- &pkg.Item{Hangul: "시험"}
		return nil
	}
	formatter := func(ctx context.Context, items <-chan *pkg.Item, w io.Writer, options map[string]string) error {
		return fmt.Errorf("boom: formatter")
	}
	options := map[string]string{}

	err := Run(in, out, parser, formatter, options)
	assert.NotNil(t, err)
	assert.Equal(t, "", out.String())
}

func TestRun_NilParser(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	formatter := func(ctx context.Context, items <-chan *pkg.Item, w io.Writer, options map[string]string) error {
		return nil
	}
	options := map[string]string{}

	err := Run(in, out, nil, formatter, options)
	assert.NotNil(t, err)
	assert.Equal(t, "", out.String())
}

func TestRun_NilFormatter(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	parser := func(ctx context.Context, r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
		items <- &pkg.Item{Hangul: "시험"}
		close(items)
		return nil
	}
	options := map[string]string{}

	err := Run(in, out, parser, nil, options)
	assert.NotNil(t, err)
	assert.Equal(t, "", out.String())
}

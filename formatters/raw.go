package formatters

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

func FormatRaw(items <-chan *pkg.Item, w io.Writer, options map[string]string) error {
	writeHeader(w, options)
	for item := range items {
		j, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(w, string(j))
	}
	return nil
}

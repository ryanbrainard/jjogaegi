package parsers

import (
	"context"
	"encoding/csv"
	"io"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

func ParseTSV(ctx context.Context, in io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	r := csv.NewReader(in)
	r.Comma = '\t'

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		item := &pkg.Item{
			NoteID:     record[0],
			ExternalID: record[1],
			Hangul:     record[2],
			Hanja:      record[3],
			Def: pkg.Translation{
				Korean:  record[4],
				English: record[5],
			},
			Pronunciation: record[6],
			AudioTag:      record[7],
			ImageTag:      record[8],
			Grade:         record[9],
			Antonym:       record[10],
			Examples: []pkg.Translation{
				{
					Korean:  record[11],
					English: record[12],
				},
				{
					Korean:  record[13],
					English: record[14],
				},
			},
		}

		items <- item
	}

	return nil
}

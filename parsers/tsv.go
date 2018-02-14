package parsers

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

func ParseTSV(in io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	r := csv.NewReader(in)
	r.Comma = '\t'

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if len(record) < 13 {
			return fmt.Errorf("Missing required fields (%d): %+v", len(record), record)
		}

		item := &pkg.Item{
			NoteID:        record[0],
			Hangul:        record[1],
			Hanja:         record[2],
			Pronunciation: record[3],
			AudioTag:      record[4],
			Def: pkg.Translation{
				Korean:  record[5],
				English: record[6],
			},
			Antonym: record[7],
			Examples: []pkg.Translation{
				{
					Korean:  record[8],
					English: record[9],
				},
				{
					Korean:  record[10],
					English: record[11],
				},
			},
			ImageTag: record[12],
		}

		items <- item
	}

	return nil
}

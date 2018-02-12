package parsers

import (
	"encoding/csv"
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

		item := &pkg.Item{
			NoteId:        record[0],
			Id:            record[1],
			Hangul:        record[2],
			Hanja:         record[3],
			Pronunciation: record[4],
			AudioTag:      record[5],
			Def: pkg.Translation{
				Korean:  record[6],
				English: record[7],
			},
			Antonym: record[8],
			Examples: []pkg.Translation{
				{
					Korean:  record[9],
					English: record[10],
				},
				{
					Korean:  record[11],
					English: record[12],
				},
			},
			ImageTag: record[13],
		}

		items <- item
	}

	return nil
}

package interceptors

import (
	"github.com/google/uuid"
	"github.com/ryanbrainard/jjogaegi/pkg"
)

func GenerateNoteId(item *pkg.Item, options map[string]string) error {
	if item.NoteID != "" {
		return nil
	}

	item.NoteID = uuid.New().String()

	return nil
}

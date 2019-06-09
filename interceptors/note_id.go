package interceptors

import (
	"github.com/google/uuid"
	"go.ryanbrainard.com/jjogaegi/pkg"
)

func GenerateNoteId(item *pkg.Item, options map[string]string) error {
	if item.NoteID == "" {
		item.NoteID = uuid.New().String()
	}

	return nil
}

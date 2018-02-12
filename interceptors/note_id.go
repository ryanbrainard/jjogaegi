package interceptors

import (
	"github.com/google/uuid"
	"github.com/ryanbrainard/jjogaegi/pkg"
)

func GenerateNoteId(item *pkg.Item, options map[string]string) error {
	if item.NoteId != "" {
		return nil
	}

	item.NoteId = uuid.New().String()

	return nil
}

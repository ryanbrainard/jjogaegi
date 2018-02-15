package interceptors

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/ryanbrainard/jjogaegi/pkg"
)

func MediaFormatting(item *pkg.Item, options map[string]string) error {
	if tag, err := formatAudioTag(item.AudioTag, options); err == nil {
		item.AudioTag = tag
	} else {
		return err
	}

	if tag, err := formatImageTag(item.ImageTag, options); err == nil {
		item.ImageTag = tag
	} else {
		return err
	}

	return nil
}

func formatAudioTag(audioTag string, options map[string]string) (string, error) {
	return formatMediaTag(audioTag, "[sound:%s]", options)
}

func formatImageTag(imageTag string, options map[string]string) (string, error) {
	return formatMediaTag(imageTag, "<img src=\"%s\">", options)
}

func formatMediaTag(mediaTag string, format string, options map[string]string) (string, error) {
	if !strings.HasPrefix(mediaTag, "http") {
		return mediaTag, nil
	}
	mediaURL := mediaTag

	mediaDirEnvKey := "ANKI_MEDIA_DIR"
	mediaDir := options[pkg.OPT_MEDIADIR]
	if mediaDir == "" {
		mediaDir = os.Getenv(mediaDirEnvKey)
	}
	if mediaDir == "" {
		panic("media dir fail safe...") // TODO: fix error bubbling and remove this
		return "", fmt.Errorf("Cannot download media (%s) if %s option or %s env is not set", mediaURL, pkg.OPT_MEDIADIR, mediaDirEnvKey)
	}

	filename := uuid.New().String() + path.Ext(mediaURL)
	err := downloadMedia(mediaURL, path.Join(mediaDir, filename))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(format, filename), nil
}

func downloadMedia(mediaURL string, filename string) error {
	resp, err := http.Get(mediaURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

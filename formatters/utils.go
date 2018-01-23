package formatters

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

func writeHeader(out io.Writer, options map[string]string) {
	if s, ok := options[pkg.OPT_HEADER]; ok && s != "" {
		out.Write([]byte(s + "\n"))
	}
}

func formatHangulHanja(item *pkg.Item, options map[string]string) string {
	switch options[pkg.OPT_HANJA] {
	case pkg.OPT_HANJA_PARENTHESIS:
		s := item.Hangul
		if len(item.Hanja) > 0 {
			s += " (" + item.Hanja + ")"
		}
		return s
	default:
		return item.Hangul
	}
}

func formatAudioTag(item *pkg.Item, options map[string]string) (string, error) {
	if item.AudioURL == "" {
		return "", nil
	}

	filename := item.AudioURL
	if strings.HasPrefix(item.AudioURL, "http") && options[pkg.OPT_MEDIADIR] != "" {
		filename = path.Base(item.AudioURL)
		err := downloadMedia(item.AudioURL, path.Join(options[pkg.OPT_MEDIADIR], filename))
		if err != nil {
			return "", err
		}
	}

	return "[sound:" + filename + "]", nil
}

// TODO: de-dupe
func formatImageTag(item *pkg.Item, options map[string]string) (string, error) {
	if item.ImageURL == "" {
		return "", nil
	}

	filename := item.ImageURL
	if strings.HasPrefix(item.ImageURL, "http") && options[pkg.OPT_MEDIADIR] != "" {
		filename = path.Base(item.ImageURL)
		err := downloadMedia(item.ImageURL, path.Join(options[pkg.OPT_MEDIADIR], filename))
		if err != nil {
			return "", err
		}
	}

	return "<img src=\"" + filename + "\">", nil
}

func downloadMedia(url string, filename string) error {
	log.Printf("download type=audio url=%q filename=%q", url, filename)

	resp, err := http.Get(url)
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

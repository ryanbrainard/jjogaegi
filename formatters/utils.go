package formatters

import (
	"github.com/ryanbrainard/jjogaegi/pkg"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
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

func formatAudioTag(item *pkg.Item, options map[string]string) string {
	if item.AudioURL == "" {
		return ""
	}

	filename := item.AudioURL
	if strings.HasPrefix(item.AudioURL, "http") && options[pkg.OPT_AUDIODIR] != "" {
		filename = path.Base(item.AudioURL)
		err := downloadAudio(item.AudioURL, path.Join(options[pkg.OPT_AUDIODIR], filename))
		if err != nil {
			log.Fatal(err)
		}
	}

	return "[sound:" + filename + "]"
}

func downloadAudio(url string, filename string) error {
	log.Printf("download url=%q filename=%q", url, filename)

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

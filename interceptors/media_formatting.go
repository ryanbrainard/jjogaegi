package interceptors

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/ryanbrainard/jjogaegi/pkg"
	"regexp"
	"bufio"
	"errors"
	"mime"
	"path"
)

func MediaFormatting(item *pkg.Item, options map[string]string) error {
	if tag, err := formatAudioTag(item.AudioTag, options); err == nil {
		item.AudioTag = tag
	} else {
		return err
	}

	if tag, err := formatImageTag(item.ImageTag, options); err == nil {
		item.ImageTag = tag
		pkg.Debug(options, "at=format.image tag=%q", item.ImageTag)
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

	mediaDir := options[pkg.OPT_MEDIADIR]
	if mediaDir == "" {
		return "", fmt.Errorf("cannot download media (%s) unless media dir is set", mediaURL)
	}

	filename, err := downloadMedia(mediaURL, mediaDir)
	if err != nil {
		return "", err
	}
	pkg.Debug(options, "at=format.media filename=%s", filename)

	return fmt.Sprintf(format, filename), nil
}

func downloadMedia(mediaURL, mediaDir string) (string, error) {
	resp, err := http.Get(mediaURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	switch contentType := strings.ToLower(resp.Header.Get("Content-Type")); contentType {
	case "text/html; charset=utf-8":
		newMediaURL := extractImageURL(resp.Body)
		if newMediaURL != "" {
			return downloadMedia(newMediaURL, mediaDir)
		}
		return "", errors.New("media is unknown HTML format: " + mediaURL)
	default:
		exts, err := mime.ExtensionsByType(contentType)
		if err != nil {
			return "", err
		}
		ext := ""
		if len(exts) > 1 {
			ext = exts[0]
		}
		filename := uuid.New().String() + ext
		filepath := path.Join(mediaDir, filename)

		file, err := os.Create(filepath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		return filename, err
	}
}

var regexpImageURL = regexp.MustCompile(`<img src="(http:\/\/dicmedia.korean.go.kr:8899\/multimedia\/multimedia_files\/convert\/.*?jpg)"`)

func extractImageURL(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if m := regexpImageURL.FindStringSubmatch(line); m != nil {
			return m[1]
		}
	}
	return ""
}
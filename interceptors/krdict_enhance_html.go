package interceptors

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"go.ryanbrainard.com/jjogaegi/pkg"
)

// TODO: temp screenscaping solution until API improves. remove this interceptor asap!
func KrDictEnhanceHTML(item *pkg.Item, options map[string]string) error {
	// format: strings.Join([]string{"krdict", lang, entryId, lexicalUnit}, ":")
	idSplit := strings.Split(item.ExternalID, ":")
	if len(idSplit) != 4 || idSplit[0] != "krdict" {
		return nil
	}
	entryID := idSplit[2]

	if item.AudioTag == "" || strings.HasPrefix(item.AudioTag, "[sound:say-") || hasSoundProblem(item, options) {
		pageReader, err := fetchHTML(entryID)
		if err != nil {
			return nil
		}
		defer pageReader.Close()

		if url := extractAudioURL(pageReader); url != "" {
			item.AudioTag = url
		}
	}

	return nil
}

func fetchHTML(entryId string) (io.ReadCloser, error) {
	url := "https://krdict.korean.go.kr/eng/dicSearch/SearchView?nation=eng&ParaWordNo=" + entryId

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

var regexpAudioURL = regexp.MustCompile("javascript:fnSoundPlay\\('(.*?mp3)'\\)")

func extractAudioURL(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if m := regexpAudioURL.FindStringSubmatch(line); m != nil {
			return m[1]
		}
	}
	return ""
}

var regexpAudioTag = regexp.MustCompile("\\[sound:(.*)\\]")

func hasSoundProblem(item *pkg.Item, options map[string]string) bool {
	m := regexpAudioTag.FindStringSubmatch(item.AudioTag)
	if m == nil {
		return false
	}

	filename := path.Join(options[pkg.OPT_MEDIADIR], m[1])

	file, err := os.Open(filename)
	if err != nil {
		return true
	}
	defer file.Close()

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return false
	}

	// Reset the read pointer if necessary.
	file.Seek(0, 0)

	contentType := http.DetectContentType(buffer)
	if strings.Contains(contentType, "text/html") {
		return true
	}

	return false
}

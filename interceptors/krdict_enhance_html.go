package interceptors

import (
	"bufio"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/ryanbrainard/jjogaegi/pkg"
)

// TODO: temp screenscaping solution until API improves. remove this interceptor asap!
func KrDictEnhanceHTML(item *pkg.Item, options map[string]string) error {
	// format: strings.Join([]string{"krdict", lang, entryId, lexicalUnit}, ":")
	idSplit := strings.Split(item.ExternalID, ":")
	if len(idSplit) != 4 || idSplit[0] != "krdict" {
		return nil
	}
	entryID := idSplit[2]

	if item.AudioTag == "" || strings.HasPrefix(item.AudioTag, "[sound:say-") {
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

func extractAudioURL(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if m := regexp.MustCompile("javascript:fnSoundPlay\\('(.*)'\\)").FindStringSubmatch(line); m != nil {
			return m[1]
		}
	}
	return ""
}

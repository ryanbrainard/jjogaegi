package parsers

import (
	"io"
	"log"
	"path"
	"strings"

	"net/http"

	"fmt"

	"os"

	"github.com/ryanbrainard/jjogaegi/pkg"
	"launchpad.net/xmlpath"
)

const supportedLanguage = "kor"
const supportedLexicalUnit = "단어"

func ParseKrDictXML(r io.Reader, items chan<- *pkg.Item, options map[string]string) error {
	rootNode, err := xmlpath.Parse(r)
	if err != nil {
		return err
	}

	lang := get(rootNode, "/LexicalResource/Lexicon/feat[@att='language']/@val")
	if lang != supportedLanguage {
		return fmt.Errorf("Only %q supported.", supportedLanguage)
	}

	entriesIter := xmlpath.MustCompile("/LexicalResource/Lexicon/LexicalEntry").Iter(rootNode)
	for {
		if !entriesIter.Next() {
			break
		}
		entryNode := entriesIter.Node()

		lexicalUnit := get(entryNode, "feat[@att='lexicalUnit']/@val")
		if lexicalUnit != supportedLexicalUnit {
			continue
		}

		entryId := get(entryNode, ".[@att='id']/@val")

		noteId := strings.Join([]string{"krdict", lang, entryId, lexicalUnit}, ":")

		audioURL := get(entryNode, "WordForm/feat[@att='sound']/@val")
		audioTag, err := formatAudioTag(audioURL, options)
		if err != nil {
			return err
		}

		imageURL := get(entryNode, "Sense/Multimedia/feat[@att='url']/@val")
		imageTag, err := formatImageTag(imageURL, options)
		if err != nil {
			return err
		}

		item := &pkg.Item{
			NoteId:        noteId,
			Id:            noteId,
			Hangul:        get(entryNode, "Lemma/feat/@val"),
			Hanja:         get(entryNode, "feat[@att='origin']/@val"),
			Pronunciation: get(entryNode, "WordForm/feat[@att='pronunciation']/@val"),
			AudioTag:      audioTag,
			ImageTag:      imageTag,
			Antonym:       get(entryNode, "Sense/SenseRelation/feat[@val='반대말']/../feat[@att='lemma']/@val"),
			Def: pkg.Translation{
				Korean: get(entryNode, "Sense/feat[@att='definition']/@val"),
				// English: def is fetcher in enhancer
			},
		}

		exampleType := ""
		examplesIter := xmlpath.MustCompile("Sense/SenseExample").Iter(entryNode)
		for {
			if !examplesIter.Next() {
				break
			}
			exampleNode := examplesIter.Node()

			nextExampleType := get(exampleNode, "feat[@att='type']/@val")
			if nextExampleType == exampleType {
				continue // choose one of each example type
			}
			if nextExampleType == "대화" {
				continue
			}
			exampleType = nextExampleType

			item.Examples = append(item.Examples, pkg.Translation{
				Korean: get(exampleNode, "feat[@att='example']/@val"),
			})
		}

		items <- item
	}

	return nil
}

func get(node *xmlpath.Node, xpath string) string {
	path := xmlpath.MustCompile(xpath)

	if value, ok := path.String(node); ok {
		return value
	}

	return ""
}

func formatAudioTag(audioURL string, options map[string]string) (string, error) {
	return formatMediaTag(audioURL, "[sound:%s]", options)
}

func formatImageTag(imageURL string, options map[string]string) (string, error) {
	return formatMediaTag(imageURL, "<img src=\"%s\">", options)
}

func formatMediaTag(mediaURL string, format string, options map[string]string) (string, error) {
	if mediaURL == "" {
		return "", nil
	}

	filename := mediaURL

	if strings.HasPrefix(mediaURL, "http") {
		if options[pkg.OPT_MEDIADIR] == "" {
			return "", fmt.Errorf("Cannot download media (%s) if %s is not set", mediaURL, pkg.OPT_MEDIADIR)
		}

		filename = path.Base(mediaURL)
		err := downloadMedia(mediaURL, path.Join(options[pkg.OPT_MEDIADIR], filename))
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf(format, filename), nil
}

func downloadMedia(mediaURL string, filename string) error {
	log.Printf("download type=audio url=%q filename=%q", mediaURL, filename)

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

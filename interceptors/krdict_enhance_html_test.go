package interceptors

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractAudioURL(t *testing.T) {
	pageFilename := "../testing/fixtures/kr_dict_en_15392.html"
	pageFile, err := os.Open(pageFilename)
	assert.NoError(t, err)
	defer pageFile.Close()

	audioURL := extractAudioURL(pageFile)
	assert.Equal(t, "http://dicmedia.korean.go.kr:8899/multimedia/multimedia_files/convert/20120306/39540/SND000028243.mp3", audioURL)
}

func TestExtractAudioURL_Multiple(t *testing.T) {
	pageFilename := "../testing/fixtures/kr_dict_en_62465.html"
	pageFile, err := os.Open(pageFilename)
	assert.NoError(t, err)
	defer pageFile.Close()

	audioURL := extractAudioURL(pageFile)
	assert.Equal(t, "http://dicmedia.korean.go.kr:8899/multimedia/multimedia_files/convert/20160913/20000/14000/318855/SND000328209.mp3", audioURL)
}

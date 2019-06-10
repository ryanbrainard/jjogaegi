package testing

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strconv"
)

func NewKrdictMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		vcrFilenameRaw := r.URL.Path
		vcrQuerty := r.URL.Query()
		vcrQuerty.Del("key")
		baseUrl := "https://krdict.korean.go.kr"
		if customBaseUrl := vcrQuerty.Get("_baseUrl"); customBaseUrl != "" {
			baseUrl = customBaseUrl
		}
		vcrFilenameRaw += "?" + vcrQuerty.Encode()

		h := md5.New()
		io.WriteString(h, vcrFilenameRaw)
		vcrFilenameHash := base64.RawURLEncoding.EncodeToString((h.Sum(nil)))

		vcrFilepathHash := "../testing/fixtures/vcr/krdict/" + vcrFilenameHash

		if os.Getenv("VCR_RECORD") == strconv.FormatBool(true) {
			baseUrl, _ := url.Parse(baseUrl)
			r.URL.Scheme = baseUrl.Scheme
			r.URL.Host = baseUrl.Host
			vcrResponse, err := http.Get(r.URL.String())
			if err != nil {
				panic(err)
			}
			defer vcrResponse.Body.Close()

			if vcrResponse.StatusCode != http.StatusOK {
				panic(vcrResponse.Status)
			}

			// TODO: copy headers?

			err = os.MkdirAll(path.Dir(vcrFilepathHash), os.ModePerm)
			if err != nil {
				panic(err)
			}

			vcrFile, err := os.Create(vcrFilepathHash)
			if err != nil {
				panic(err)
			}
			defer vcrFile.Close()

			_, err = io.Copy(vcrFile, vcrResponse.Body)
			if err != nil {
				panic(err)
			}

			vcrFile.Close()
		}

		vcrFile, err := os.Open(vcrFilepathHash)
		if err == nil {
			defer vcrFile.Close()
			vcrData, err := ioutil.ReadAll(vcrFile)
			if err != nil {
				panic(err)
			}

			w.Header().Set("Content-Type", http.DetectContentType(vcrData))
			_, err = io.Copy(w, bytes.NewBuffer(vcrData))
			if err != nil {
				panic(err)
			}

			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Unknown path " + r.URL.Path + ". Run with VCR_RECORD=true to record real responses"))
		return
	}))
}

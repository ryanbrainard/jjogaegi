package testing

import (
	"io"
	"net/http"
	"net/http/httptest"
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

		vcrFilename := "../testing/fixtures/vcr/krdict/" + r.URL.Path
		vcrQuerty := r.URL.Query()
		vcrQuerty.Del("key")
		vcrFilename += "?" + vcrQuerty.Encode()

		if os.Getenv("VCR_RECORD") == strconv.FormatBool(true) {
			r.URL.Scheme = "https"
			r.URL.Host = "krdict.korean.go.kr"
			vcrResponse, err := http.Get(r.URL.String())
			if err != nil {
				panic(err)
			}
			defer vcrResponse.Body.Close()

			if vcrResponse.StatusCode != http.StatusOK {
				panic(vcrResponse.Status)
			}

			// TODO: copy headers?

			err = os.MkdirAll(path.Dir(vcrFilename), os.ModePerm)
			if err != nil {
				panic(err)
			}

			vcrFile, err := os.Create(vcrFilename)
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

		vcrFile, err := os.Open(vcrFilename)
		if err == nil {
			defer vcrFile.Close()
			_, err = io.Copy(w, vcrFile)
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

package testing

import (
	"os"
	"net/http"
	"net/http/httptest"
	"io"
)

func NewKrdictMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var filename string

		switch r.URL.Path {
		case "/api/view":
			// TODO: test other query params
			filename = "kr_dict_en_" + r.URL.Query().Get("q") + ".xml"
		case "/api/search":
			// TODO: test other query params
			filename = "kr_dict_en_search_"  + r.URL.Query().Get("q") + ".xml"
		case "/multimedia/multimedia_files/convert/20150929/201390/dummy.jpg":
			// TODO: more generic?
			filename = "dummy.jpg"
		}

		if filename == "" {
			w.WriteHeader(404)
			w.Write([]byte("Unknown path " + r.URL.Path))
			return
		}

		// TODO: make path more resilient
		fixture, err := os.Open("../testing/fixtures/" + filename)
		if err != nil {
			panic(err)
		}

		if _, err = io.Copy(w, fixture); err != nil {
			panic(err)
		}
	}))
}

package main


import (
	"fmt"
	"net/http"
	"strings"
)

type MockRequester interface {
	MockGet(url string)string
}

type UrlGetServer struct {
	page MockRequester
}

func (u *UrlGetServer)ServeHTTP(w http.ResponseWriter, r *http.Request){
	url := strings.TrimPrefix(r.URL.Path, "/yandex/")
	urlcrw := u.page.MockGet(url)

	if urlcrw == "nil" {
		w.WriteHeader((http.StatusNotFound))
	}
	fmt.Fprintf(w, urlcrw)
}
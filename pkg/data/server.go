package data

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type server struct {
	fetch *client
}

type fetchRequest struct {
	Code string `json:"code"`
}

func NewServer(fetchClient *client) server {
	return server{
		fetch: fetchClient,
	}
}

func (s server) Fetch(w http.ResponseWriter, req *http.Request) {
	log.Println("handling fetch")
	path := req.URL.Path
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		serviceError(w, err)
	}
	var fetchReq fetchRequest
	err = json.Unmarshal(b, &fetchReq)

	resp, err := s.fetch.Fetch(path, fetchReq.Code)
	if err != nil {
		serviceError(w, err)
	}

	w.Write(resp)
}

func serviceError(w http.ResponseWriter, err error) {
	log.Printf("service error: %s\n", err.Error())
	w.WriteHeader(500)
	w.Write([]byte("error: " + err.Error()))
}

func (s server) DefineServer(mux *http.ServeMux) {
	mux.HandleFunc("/", withPathParam("/fetch/", s.Fetch))
}

func withPathParam(prefix string, fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.URL.Path, prefix) {
			return
		}

		path := strings.TrimPrefix(req.URL.Path, prefix)
		log.Printf("trimmed %s -> %s", req.URL.Path, path)
		req.URL.Path = path
		fn(w, req)
	}
}

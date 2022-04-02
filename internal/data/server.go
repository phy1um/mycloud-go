package data

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type server struct {
	fetch *client
}

type fetchRequest struct {
	Key  string `json:"key"`
	Code string `json:"code"`
}

func NewServer(fetchClient *client) server {
	return server{
		fetch: fetchClient,
	}
}

func (s server) Fetch(w http.ResponseWriter, req *http.Request) {
	log.Println("handling fetch")
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		serviceError(w, err)
	}
	var fetchReq fetchRequest
	err = json.Unmarshal(b, &fetchReq)

	resp, err := s.fetch.Fetch(fetchReq.Key, fetchReq.Code)
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
	mux.HandleFunc("/fetch", s.Fetch)
}

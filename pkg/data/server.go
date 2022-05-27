package data

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type server struct {
	fetch *client
	ctx   context.Context
}

type fetchRequest struct {
	Code string `json:"code"`
}

func NewServer(ctx context.Context, fetchClient *client) server {
	return server{
		ctx:   ctx,
		fetch: fetchClient,
	}
}

func (s server) Fetch(w http.ResponseWriter, req *http.Request) {
	ctx := s.ctx
	log.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("request-path", req.URL.Path).
			Str("request-method", req.Method)
	})

	path := req.URL.Path
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		serviceError(ctx, w, err)
	}
	var fetchReq fetchRequest
	err = json.Unmarshal(b, &fetchReq)

	resp, err := s.fetch.Fetch(ctx, path, fetchReq.Code)
	if err != nil {
		serviceError(ctx, w, err)
	}

	w.Write(resp)
}

func serviceError(ctx context.Context, w http.ResponseWriter, err error) {
	log.Ctx(ctx).Err(err).Msg("service error")
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
		//log.Printf("trimmed %s -> %s", req.URL.Path, path)
		req.URL.Path = path
		fn(w, req)
	}
}

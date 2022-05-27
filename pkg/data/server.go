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

type ctxHttpFunc func(context.Context, http.ResponseWriter, *http.Request)

func NewServer(ctx context.Context, fetchClient *client) server {
	return server{
		ctx:   ctx,
		fetch: fetchClient,
	}
}

func (s server) Fetch(ctx context.Context, w http.ResponseWriter, req *http.Request) {
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
	log.Ctx(ctx).Error().Stack().Err(err).Msg("service error")
	w.WriteHeader(500)
	w.Write([]byte("error: " + err.Error()))
}

func (s server) DefineServer(mux *http.ServeMux) {
	mux.HandleFunc("/", withPathParam(s.ctx, "/fetch/", s.Fetch))
}

func withPathParam(ctx context.Context, prefix string, fn ctxHttpFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.URL.Path, prefix) {
			return
		}

		wlog := log.Ctx(ctx).With().
			Str("http-path", req.URL.Path).
			Str("http-method", req.Method).
			Dict("http-headers", zerolog.Dict().Fields(req.Header)).
			Logger()

		wctx := wlog.WithContext(ctx)

		path := strings.TrimPrefix(req.URL.Path, prefix)
		req.URL.Path = path

		fn(wctx, w, req)
	}
}

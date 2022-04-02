package internal

import (
	"fmt"
	"net/http"
)

type Health struct {
	Version string
}

func (h Health) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	msg := fmt.Sprintf("running %s", h.Version)
	w.Write([]byte(msg))
}

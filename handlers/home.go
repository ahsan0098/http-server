package handlers

import (
	"net/http"
)

type Home struct{}

func HomeHdlr() *Home {
	return &Home{}
}

func (h *Home) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("welcome to home"))
}

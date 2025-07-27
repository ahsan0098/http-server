package handlers

import (
	"corenethttp/helpers/jsons"
	"log"
	"net/http"
)

type Hello struct {
	lg *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h.lg.Println("Hello World")

	d := r.URL.Query()

	jsons.JsonResponse(w, http.StatusOK, d)
}

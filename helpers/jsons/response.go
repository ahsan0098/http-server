package jsons

import (
	"encoding/json"
	"log"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}

func ErrorResponse(w http.ResponseWriter, code int, msg string) {

	if code > 499 {
		log.Printf("Internal %v Error : %v \n", code, msg)
		w.WriteHeader(500)
		return
	}

	type errREsponse struct {
		Error string `json:"error"`
	}

	JsonResponse(w, code, errREsponse{Error: msg})
}

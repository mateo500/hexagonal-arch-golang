package http

import (
	"log"
	"net/http"
)

func SetupResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Printf("error seting up http response: %v", err)
	}
}

package service

import (
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(`{"status":"running"}`))
}

func (s *DCAService) StartWebServer() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":80", nil)
}

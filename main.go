package main

import (
	"net/http"
)

func main() {
	const fileRootPath = "."
	const port = "8080"
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(fileRootPath))))
	serveMux.HandleFunc("/healthz", readinessHandler)
	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}
	server.ListenAndServe()
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
	r.Header.Set("Content-Type", "text/plain; charset=utf-8")
}

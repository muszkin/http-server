package main

import (
	"net/http"
)

func main() {
	const fileRootPath = "."
	const port = "8080"
	serveMux := http.NewServeMux()	
	serveMux.Handle("/", http.FileServer(http.Dir(fileRootPath)))

	server := http.Server{
		Handler: serveMux,
		Addr: ":" + port,
	}
	server.ListenAndServe()
}

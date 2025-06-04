package main

import (
	"net/http"
)


func main() {
	mux := http.NewServeMux()

	server := http.Server{
		Handler:	mux,
		Addr:		":8080",
	}

	fileServerHandler := http.FileServer(http.Dir("."))
	mux.Handle("/app/", http.StripPrefix("/app/", fileServerHandler))
	mux.HandleFunc("/", readinessHandler)


	server.ListenAndServe()
}

func readinessHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}

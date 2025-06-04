package main

import (
	"net/http"
	"sync/atomic"
	"fmt"
)


type apiConfig struct {
	fileserverHits atomic.Int32
}

type chirpBody struct {
	Body string
}

func main() {
	mux := http.NewServeMux()

	server := http.Server{
		Handler:	mux,
		Addr:		":8080",
	}

	apiConfig := apiConfig{}

	mux.HandleFunc("GET /admin/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiConfig.getRequestCountsHandler)
	mux.HandleFunc("POST /admin/reset", apiConfig.resetRequestCountsHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	fileServerHandler := http.FileServer(http.Dir("."))
	fileServerHandlerNoPrefix := http.StripPrefix("/app/", fileServerHandler)
	mux.Handle("/app/", apiConfig.middlewareMetricsInc(fileServerHandlerNoPrefix))


	server.ListenAndServe()
}


func readinessHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}

func (cfg *apiConfig) getRequestCountsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	htmlString :=`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	outputString := fmt.Sprintf(htmlString, cfg.fileserverHits.Load())
	w.Write([]byte(outputString))
}

func (cfg *apiConfig) resetRequestCountsHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Request Counts Reset"))
}



func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

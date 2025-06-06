package main

import (
	"net/http"
	"github.com/wexlerdev/chirpy/internal/handlers"
)


import _ "github.com/lib/pq"


type chirpBody struct {
	Body string
}

func main() {

	mux := http.NewServeMux()

	server := http.Server{
		Handler:	mux,
		Addr:		":8080",
	}

	api := handlers.NewAPI()

	mux.HandleFunc("GET /admin/healthz", api.ReadinessHandler)
	mux.HandleFunc("GET /admin/metrics", api.GetRequestCountsHandler)
	//reset endpoint
	resetHandlerFunc := http.HandlerFunc(api.ResetUsersHandler)
	wrappedDevOnlyHandler := api.DevOnlyMiddleware(resetHandlerFunc)
	mux.Handle("POST /admin/reset", wrappedDevOnlyHandler)

	mux.HandleFunc("POST /api/validate_chirp", api.ValidateChirpHandler)
	mux.HandleFunc("POST /api/users", api.CreateUserHandler)

	fileServerHandler := http.FileServer(http.Dir("."))
	fileServerHandlerNoPrefix := http.StripPrefix("/app/", fileServerHandler)
	mux.Handle("/app/", api.MiddlewareMetricsInc(fileServerHandlerNoPrefix))


	server.ListenAndServe()
}




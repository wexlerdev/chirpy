package main

import (
	"github.com/wexlerdev/chirpy/internal/handlers"
	"net/http"
	"fmt"
)

import _ "github.com/lib/pq"

type chirpBody struct {
	Body string
}

func main() {

	mux := http.NewServeMux()

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	api := handlers.NewAPI()

	mux.HandleFunc("GET /admin/healthz", api.ReadinessHandler)
	mux.HandleFunc("GET /admin/metrics", api.GetRequestCountsHandler)
	//reset endpoint
	resetHandlerFunc := http.HandlerFunc(api.ResetUsersHandler)
	wrappedDevOnlyHandler := api.DevOnlyMiddleware(resetHandlerFunc)
	mux.Handle("POST /admin/reset", wrappedDevOnlyHandler)

	mux.HandleFunc("POST /api/chirps", api.HandleCreateChirp)
	mux.HandleFunc("GET /api/chirps/{id}", api.HandleGetChirp)
	mux.HandleFunc("GET /api/chirps", api.HandleGetAllChirps)
	mux.HandleFunc("POST /api/users", api.CreateUserHandler)
	mux.HandleFunc("POST /api/login", api.HandleLogin)
	mux.HandleFunc("POST /api/refresh", api.HandleRefresh)
	mux.HandleFunc("POST /api/revoke", api.HandleRevoke)
	mux.HandleFunc("PUT /api/users", api.HandleUpdateUser)
	mux.HandleFunc("DELETE /api/chirps/{id}", api.HandleDeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", api.HandleUserUpgrade)

	fileServerHandler := http.FileServer(http.Dir("."))
	fileServerHandlerNoPrefix := http.StripPrefix("/app/", fileServerHandler)
	mux.Handle("/app/", api.MiddlewareMetricsInc(fileServerHandlerNoPrefix))

	fmt.Println("about to listen and serve")

	server.ListenAndServe()
}

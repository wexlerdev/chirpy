package handlers

import (
	"github.com/wexlerdev/chirpy/internal/config"
	"fmt"
	"net/http"
)


type API struct {
	cfg * config.ApiConfig
}

func NewAPI() *API {
	apiConfig := config.NewConfig()

	api := API{
		cfg:	apiConfig,
	}
	return &api
}


func (api *API) GetRequestCountsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	htmlString :=`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	outputString := fmt.Sprintf(htmlString, api.cfg.FileserverHits.Load())
	w.Write([]byte(outputString))
}


func (api *API) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (api *API) ReadinessHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("OK"))
}

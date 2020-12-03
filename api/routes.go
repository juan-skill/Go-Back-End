package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/other_project/crockroach/api/httphand"
	db "github.com/other_project/crockroach/internal/storage"
)

// Routes create an router multiplexer
func Routes(store *db.Store) *chi.Mux {
	mux := chi.NewMux()

	// globals middleware
	mux.Use(
		middleware.Logger,    //log every http request
		middleware.Recoverer, // recover if a panic occurs
	)

	handler := httphand.NewHandlerRequest(store)

	mux.Get("/status", showStatus)
	mux.Post("/domain", handler.Create)

	return mux
}

// showStatus return the status of the API
func showStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("done-by", "juan")

	res := map[string]interface{}{"status": "OK"}

	_ = json.NewEncoder(w).Encode(res)
}

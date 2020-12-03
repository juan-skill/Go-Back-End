package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"github.com/other_project/crockroach/internal/storage"
	"github.com/other_project/crockroach/shared/cockroachdb"
	"github.com/other_project/crockroach/shared/env"
)

const (
	// ReadTimeout ...
	ReadTimeout = 10 * time.Second
	// WriteTimeout ...
	WriteTimeout = 10 * time.Second
)

var (
	// PortServer data to connect with http server
	PortServer = env.GetString("SERVER_PORT", ":8090")
)

// MyServer serves HTTP requests for our service.
type MyServer struct {
	server   *http.Server
	router   *chi.Mux
	clientDB *sql.DB
}

// NewServer initialize the server instance
func NewServer(mux *chi.Mux) *MyServer {
	s := &http.Server{
		Addr:           PortServer,
		Handler:        mux,
		ReadTimeout:    ReadTimeout,
		WriteTimeout:   WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	storage.CockroachClient = *cockroachdb.NewSQLClient()

	myServer := new(MyServer)

	myServer.server = s
	myServer.router = mux
	myServer.clientDB = &storage.CockroachClient

	return myServer
}

// Run launch the  server
func (s *MyServer) Run() {
	log.Fatal(s.server.ListenAndServe())
}

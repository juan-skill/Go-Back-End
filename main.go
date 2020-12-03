package main

import (
	"github.com/other_project/crockroach/api"
	"github.com/other_project/crockroach/internal/logs"
	"github.com/other_project/crockroach/internal/storage"
)

func main() {
	_ = logs.InitLogger()

	store := storage.NewStore()
	mux := api.Routes(store)
	server := api.NewServer(mux)
	server.Run()
}

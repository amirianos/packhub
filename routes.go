package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"packhub/handlers"
)

func getRoutes(h *handlers.Repository) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", h.Home)

	return mux
}

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func Routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{"https://*", "http://*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSFR-Token"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: true,
				MaxAge:           300,
			}))

	return mux
}

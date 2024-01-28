package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

var handlerOptions = cors.Options{
	AllowedOrigins:   []string{"https://*", "http://*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSFR-Token"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
}

func (app *Config) Routes() http.Handler {
	fmt.Println("inja3")
	mux := chi.NewRouter()
	mux.Use(
		cors.Handler(handlerOptions),
		middleware.Heartbeat("/ping"),
	)

	mux.Post("/", app.Broker)

	return mux
}
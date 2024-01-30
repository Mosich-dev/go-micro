package main

import (
	"database/sql"
	"fmt"
	"github.com/Mosich-dev/go-micro/authentication-service/data"
	"log"
	"net/http"
)

const port = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting the Authentication service.")

	app := Config{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

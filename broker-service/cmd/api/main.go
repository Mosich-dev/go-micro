package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "80"

type Config struct{}

func main() {
	app := Config{}
	log.Printf("Starting the Broker Service on Port :%s", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

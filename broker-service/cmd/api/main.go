package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "80"

func main() {
	log.Printf("Starting the Broker Service on Port :%s", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: Routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

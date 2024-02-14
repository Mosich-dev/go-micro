package main

import (
	"context"
	"fmt"
	"github.com/Mosich-dev/go-micro/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

const (
	webPort            = "80"
	dbName             = "logs"
	mongoURI           = "mongodb://localhost:27017"
	mongoConnectTimout = time.Second * 45
)

var mongoClient *mongo.Client
var err error

type Config struct {
	Models data.Models
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), mongoConnectTimout)
	defer cancel()
	mongoClient, err = connectToMongo(ctx)
	if err != nil {
		log.Println("Connection to Mongo failed.")
		log.Panic(err)
	}
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	app := Config{
		Models: *data.New(mongoClient, dbName),
	}
	log.Println("Starting logger-service on port", webPort)
	// go app.Serve()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic("Failed to Listen and Serve.")
	}
}

func (app *Config) Serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic("Failed to Listen and Serve.")
	}

}

func connectToMongo(ctx context.Context) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return client, nil
}

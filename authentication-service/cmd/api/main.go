package main

import (
	"database/sql"
	"fmt"
	"github.com/Mosich-dev/go-micro/authentication-service/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
	"time"
)

const port = "80"
const connCheckTime = 2

var count int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting the Authentication service.")

	conn := connectToDB()
	if conn == nil {
		log.Panic("database connection failed.")
	}

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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("database not ready yet...")
			count++
		} else {
			log.Printf("Connected to Database at: %s\n", dsn)
			return connection
		}

		if count > 10 {
			log.Println(err)
			return nil
		}
		log.Printf("Checking the database connection every %d second(s)...", connCheckTime)
		time.Sleep(connCheckTime * time.Second)
		continue
	}
}

package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// connect to database
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to postgres")
	}
	// set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn)}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
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
	var count int
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres is not yet ready ...")
			count++
		} else {
			log.Println("Connected to postgres")
			return connection
		}

		if count >= 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for two seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}

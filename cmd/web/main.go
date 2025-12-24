package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/sillkiw/url-shorten/internal/app"

	_ "github.com/lib/pq"
)

func main() {

	connStr := "user=sillkiw password=123qweasdzxc dbname=urldb sslmode=disable"

	store, err := connectDB(connStr)
	if err != nil {
		log.Fatalf("failed to connect to db %s", err.Error())
	}
	defer store.Close()
	log.Printf("connect to db")

	app := app.New(store)

	err = http.ListenAndServe(":4000", app.Routes())
	log.Fatal(err)
}

func connectDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect")
	}
	return db, nil

}

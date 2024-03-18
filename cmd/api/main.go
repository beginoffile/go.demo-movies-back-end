package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

const port = 8080

type application struct {
	DNS    string
	Domain string
	DB     *pgx.Conn
}

func main() {
	//setup application config
	var app application

	//read from command line
	flag.StringVar(&app.DNS, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.Parse()
	//connect to the database
	app.Domain = "example.com"
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.DB = conn
	defer app.DB.Close(context.Background())

	log.Println("Starting application on port ", port)

	// http.HandleFunc("/", Hello)

	//start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}

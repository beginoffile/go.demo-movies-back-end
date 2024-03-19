package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
)

const port = 8080

type application struct {
	DNS    string
	Domain string
	DB     repository.DatabaseRepo
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

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	defer app.DB.Connection().Close(context.Background())

	log.Println("Starting application on port ", port)

	//start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}

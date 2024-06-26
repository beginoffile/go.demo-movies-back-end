package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const port = 8080

type application struct {
	DNS          string
	Domain       string
	DB           repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	APIKey       string
}

func main() {
	//setup application config
	var app application

	//read from command line
	flag.StringVar(&app.DNS, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection string")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "signed secret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "signed issuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "signed audience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.StringVar(&app.APIKey, "api-key", "b37ea4d3c8f91a72af0e1ea1d296a505", "api key")

	flag.Parse()
	//connect to the database

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	defer app.DB.Connection().Close(context.Background())

	app.auth = Auth{
		Issuer:        app.JWTIssuer,
		Audience:      app.JWTAudience,
		Secret:        app.JWTSecret,
		TokenExpiry:   time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:    "/",
		CookieName:    "refresh_token",
		CookieDomain:  app.CookieDomain,
	}

	log.Println("Starting application on port ", port)

	//start a web server
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	//http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}

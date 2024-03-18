package main

import (
	"context"
	"log"

	_ "github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func openDB(dsn string) (*pgx.Conn, error) {
	// db, err := sql.Open("pgx", dsn)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return conn, nil

}

func (app *application) connectToDB() (*pgx.Conn, error) {
	connection, err := openDB(app.DNS)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Postgress!")
	return connection, nil
}

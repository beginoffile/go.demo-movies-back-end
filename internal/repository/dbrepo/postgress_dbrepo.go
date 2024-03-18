package dbrepo

import (
	"backend/internal/models"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type PostgresDBRepo struct {
	DB *pgx.Conn
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) AllMovies() ([]*models.Movie, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	var movies []*models.Movie

	return movies, nil

}

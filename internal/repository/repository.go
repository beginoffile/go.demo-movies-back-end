package repository

import (
	"backend/internal/models"

	"github.com/jackc/pgx/v5"
)

type DatabaseRepo interface {
	Connection() *pgx.Conn
	AllMovies() ([]*models.Movie, error)
}

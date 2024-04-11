package repository

import (
	"backend/internal/models"

	"github.com/jackc/pgx/v5"
)

type DatabaseRepo interface {
	Connection() *pgx.Conn
	AllMovies(genre ...int) ([]*models.Movie, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	OneMovie(id int) (*models.Movie, error)
	OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error)
	AllGenres() ([]*models.Genre, error)
	InsertMovie(movie models.Movie) (int, error)
	UpdateMovieGenres(id int, genReIDs []int) error
	UpdateMovie(movie models.Movie) error
	DeleteMovie(id int) error
}

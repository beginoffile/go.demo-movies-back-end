package dbrepo

import (
	"backend/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type PostgresDBRepo struct {
	DB *pgx.Conn
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *pgx.Conn {
	return m.DB
}

func (m *PostgresDBRepo) AllMovies() ([]*models.Movie, error) {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `
	Select id, title, release_date, runtime, mpaa_rating, description, coalesce(image,''), created_at, updated_at
	From movies t1
	Order by t1.title
		`
	rows, err := m.DB.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movies []*models.Movie

	for rows.Next() {
		var movie models.Movie

		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.MPAARating,
			&movie.Description,
			&movie.Image,
			&movie.CreateAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		fmt.Println(movie)

		movies = append(movies, &movie)
	}

	return movies, nil

}

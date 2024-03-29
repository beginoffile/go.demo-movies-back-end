package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
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
			&movie.Runtime,
			&movie.MPAARating,
			&movie.Description,
			&movie.Image,
			&movie.CreateAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	return movies, nil

}

func (m *PostgresDBRepo) OneMovie(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `select id, title, release_date, runtime, mpaa_rating, description, coalesce(image,''), create_at, updated_at	
	from movies
	Where id = $1`

	row := m.DB.QueryRow(ctx, query, id)

	var movie models.Movie

	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreateAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	//get genres, if any
	query = `select t2.id, t2.genre
	from movies_genres t1
		Inner Join genres t2
		   On t2.id = t1.genre_id
	Where t1.movie_id = $1`

	rows, err := m.DB.Query(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	defer rows.Close()

	var genres []*models.Genre
	for rows.Next() {

		var g models.Genre

		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &g)
	}

	movie.Genres = genres

	return &movie, nil

}

func (m *PostgresDBRepo) OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `select id, title, release_date, runtime, mpaa_rating, description, coalesce(image,''), create_at, updated_at	
	from movies
	Where id = $1`

	row := m.DB.QueryRow(ctx, query, id)

	var movie models.Movie

	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreateAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, nil, err
	}

	//get genres, if any
	query = `select t2.id, t2.genre
	from movies_genres t1
		Inner Join genres t2
		   On t2.id = t1.genre_id
	Where t1.movie_id = $1`

	rows, err := m.DB.Query(ctx, query, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}

	defer rows.Close()

	var genres []*models.Genre
	var genresArray []int
	for rows.Next() {

		var g models.Genre

		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, nil, err
		}
		genres = append(genres, &g)
		genresArray = append(genresArray, g.ID)
	}

	movie.Genres = genres
	movie.GenresArray = genresArray

	var allGenres []*models.Genre

	query = `select t1.id, t1.genre
	from genres t1		
	Order by genre`

	rows, err = m.DB.Query(ctx, query)

	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, nil, err
		}
		allGenres = append(allGenres, &g)
	}

	return &movie, allGenres, nil

}

func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `select id, email, first_name, last_name, password, created_at, updated_at
	from users
	Where email = $1`

	var user models.User

	row := m.DB.QueryRow(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreateAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `select id, email, first_name, last_name, password, created_at, updated_at
	from users
	Where id = $1`

	var user models.User

	row := m.DB.QueryRow(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreateAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

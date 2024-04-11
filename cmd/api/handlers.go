package main

import (
	"backend/internal/graph"
	"backend/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go Movies up and running",
		Version: "1.0.0",
	}

	// out, err := json.Marshal(payload)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(out)
	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	// var movies []models.Movie

	// rd, _ := time.Parse("2006-01-02", "1986-03-07")

	// highlander := models.Movie{
	// 	ID:          1,
	// 	Title:       "HaightLander",
	// 	ReleaseDate: rd,
	// 	MPAARating:  "R",
	// 	Runtime:     116,
	// 	Description: "A Very nice movie",
	// 	CreateAt:    time.Now(),
	// 	UpdatedAt:   time.Now(),
	// }

	// movies = append(movies, highlander)

	// rd, _ = time.Parse("2006-01-02", "1981-06-12")

	// rota := models.Movie{
	// 	ID:          2,
	// 	Title:       "Raiders of the last Ark",
	// 	ReleaseDate: rd,
	// 	MPAARating:  "PG-13",
	// 	Runtime:     115,
	// 	Description: "Another A Very nice movie",
	// 	CreateAt:    time.Now(),
	// 	UpdatedAt:   time.Now(),
	// }

	// movies = append(movies, rota)

	movies, err := app.DB.AllMovies()
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	// out, err := json.Marshal(movies)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(out)

	_ = app.writeJSON(w, http.StatusOK, movies)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read  json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate user against  database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		_ = app.errorJSON(w, errors.New("invalid Credentials"), http.StatusBadRequest)
		return
	}

	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		_ = app.errorJSON(w, errors.New("INVALID CREDENTIALS"), http.StatusBadRequest)
		return
	}

	// create a jwt user

	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	// Generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)

	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			//parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			//get the user id from the token claims
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.errorJSON(w, errors.New("unknow user"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				app.errorJSON(w, errors.New("unknow user"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			refreshCookie := app.auth.GetRefreshCookie(tokenPairs.RefreshToken)
			http.SetCookie(w, refreshCookie)

			app.writeJSON(w, http.StatusOK, tokenPairs)

		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpriredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) MovieCatalog(w http.ResponseWriter, r *http.Request) {

	movies, err := app.DB.AllMovies()
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)

}

func (app *application) GetMovie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, err := app.DB.OneMovie(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movie)

}

func (app *application) MovieForEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, genres, err := app.DB.OneMovieForEdit(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload = struct {
		Movie  *models.Movie   `json:"movie"`
		Genres []*models.Genre `json:"genres"`
	}{
		movie,
		genres,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *application) AllGenres(w http.ResponseWriter, r *http.Request) {

	genres, err := app.DB.AllGenres()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, genres)

}

func (app *application) InsertMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie

	err := app.readJSON(w, r, &movie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// try to get an image
	movie = app.getPoster(movie)

	movie.CreateAt = time.Now()
	movie.UpdatedAt = time.Now()

	// Insert the movie
	newID, err := app.DB.InsertMovie(movie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// now handles genres

	err = app.DB.UpdateMovieGenres(newID, movie.GenresArray)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONReposponse{
		Error:   false,
		Message: "movie updated",
	}

	app.writeJSON(w, http.StatusAccepted, resp)

}

func (app *application) getPoster(movie models.Movie) models.Movie {
	type TheMovieDB struct {
		Page    int `json:"page"`
		Results []struct {
			PosterPath string `json:"poster_path"`
		} `json:"results"`
		TotalPages int `json:"total_pages"`
	}

	client := &http.Client{}
	theUrl := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s", app.APIKey)

	req, err := http.NewRequest("GET", theUrl+"&query="+url.QueryEscape(movie.Title), nil)

	if err != nil {
		log.Println(err)
		return movie
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return movie
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return movie
	}

	var responseObject TheMovieDB

	json.Unmarshal(bodyBytes, &responseObject)

	if len(responseObject.Results) > 0 {
		movie.Image = responseObject.Results[0].PosterPath
	}

	return movie

}

func (app *application) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	var payload models.Movie

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	ptrMovie, err := app.DB.OneMovie(payload.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	ptrMovie.Title = payload.Title
	ptrMovie.Description = payload.Description
	ptrMovie.ReleaseDate = payload.ReleaseDate
	ptrMovie.MPAARating = payload.MPAARating
	ptrMovie.Runtime = payload.Runtime
	ptrMovie.UpdatedAt = payload.UpdatedAt

	err = app.DB.UpdateMovie(*ptrMovie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.UpdateMovieGenres(ptrMovie.ID, payload.GenresArray)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONReposponse{
		Error:   false,
		Message: "movie updated",
	}

	app.writeJSON(w, http.StatusAccepted, resp)

}

func (app *application) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.DeleteMovie(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONReposponse{
		Error:   false,
		Message: "movie deleted",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) allMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	movies, err := app.DB.AllMovies(id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, movies)

}

func (app *application) MoviesGraphQL(w http.ResponseWriter, r *http.Request) {
	// populate our Graphql  type with the movies
	movies, _ := app.DB.AllMovies()

	// get the query from the rest
	q, _ := io.ReadAll(r.Body)
	query := string(q)

	// create a new variable of type *graphq.Graph
	g := graph.New(movies)

	// set the query string on the variable
	g.QueryString = query

	// perform the query
	resp, err := g.Query()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// send the response
	j, _ := json.MarshalIndent(resp, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

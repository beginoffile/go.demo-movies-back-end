package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	// create a router mux

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	// public routes
	mux.Get("/", app.Home)

	mux.Post("/authenticate", app.authenticate)

	mux.Get("/refresh", app.refreshToken)

	mux.Get("/logout", app.logout)

	mux.Get("/movies", app.AllMovies)

	// private routes
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.authRequired)

		mux.Get("/movies", app.MovieCatalog)

	})

	return mux
}
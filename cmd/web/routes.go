package main

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Timeout(60 * time.Second))

	mux.Route("/", func(r chi.Router) {
		r.Get("/", app.home)

		r.Route("/api", func(r chi.Router) {
			r.Route("/v1", func(r chi.Router) {
				r.Route("/users", func(r chi.Router) {
					r.Get("/", app.showUsers)
					r.Post("/", app.createUser)

					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", app.getUser)
						r.Patch("/", app.updateUser)
						r.Delete("/", app.deleteUser)
					})
				})
			})
		})
	})

	return mux
}

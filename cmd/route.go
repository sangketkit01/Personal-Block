package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sangketkit01/personal-block/internal/handlers"
	"net/http"
)

func route() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/login", handlers.Repo.LoginPage)
	mux.Post("/login/verify",handlers.Repo.LoginVerify)
	mux.Get("/signup", handlers.Repo.SignUpPage)
	mux.Post("/signup/insert", handlers.Repo.SignUpInsert)

	mux.Group(func (r chi.Router)  {
		r.Use(Auth)

		r.Post("/new-post",handlers.Repo.NewPost)
		r.Get("/",handlers.Repo.Home)
	})

	fileServer := http.FileServer(http.Dir("../../static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

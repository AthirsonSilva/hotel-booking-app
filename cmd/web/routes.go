package main

import (
	"net/http"

	"github.com/AthirsonSilva/golang-net-http-restapi/pkg/config"
	"github.com/AthirsonSilva/golang-net-http-restapi/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(NoSurf)
	router.Use(SessionLoad)
	router.Use(WriteToConsole)

	router.Get("/", handlers.Repo.Home)
	router.Get("/about", handlers.Repo.About)

	return router
}
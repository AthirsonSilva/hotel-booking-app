package main

import (
	"net/http"

	"github.com/AthirsonSilva/golang-net-http-restapi/internal/config"
	"github.com/AthirsonSilva/golang-net-http-restapi/internal/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Creates a new router to distribute all available endpoints
func routes(app *config.AppConfig) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(NoSurf)
	router.Use(SessionLoad)
	router.Use(WriteToConsole)

	router.Get("/", usecases.Repo.Home)
	router.Get("/about", usecases.Repo.About)
	router.Get("/generals-quarters", usecases.Repo.Generals)
	router.Get("/majors-suite", usecases.Repo.Majors)

	router.Get("/reservation-summary", usecases.Repo.ReservationSummary)
	router.Get("/make-reservation", usecases.Repo.MakeReservation)
	router.Post("/make-reservation", usecases.Repo.PostReservation)

	router.Get("/search-availability", usecases.Repo.Availability)
	router.Post("/search-availability-json", usecases.Repo.PostAvailabilityJSON)
	router.Post("/search-availability", usecases.Repo.PostAvailability)
	router.Get("/choose-room/{id}", usecases.Repo.ChooseRoom)

	router.Get("/login", usecases.Repo.LoginPage)
	router.Post("/login", usecases.Repo.Login)

	router.Get("/logout", usecases.Repo.Logout)

	fileServer := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	router.Route("/admin", func(router chi.Router) {
		// router.Use(VerifyUserAuthentication)
		router.Get("/dashboard", usecases.Repo.AdminDashboard)

		router.Get("/reservations/new", usecases.Repo.AdminAllNewReservations)
		router.Get("/reservations/all", usecases.Repo.AdminAllReservations)
		router.Get("/reservation/calendar", usecases.Repo.AdminReservationsCalendar)
		router.Get("/reservations/{id}", usecases.Repo.AdminShowSingleReservation)
	})

	return router
}

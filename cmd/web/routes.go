package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/goldalee/golangprojects/bookings/internal/config"
	"github.com/goldalee/golangprojects/bookings/internal/handlers"
)

//I had to install pat router by typing in: go get github.com/bmizerany/pat
func routes(app *config.AppConfig) http.Handler {
	//this is using pat
	// mux := pat.New()
	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	//return mux

	//this is using chi
	mux := chi.NewRouter()

	//Using middleware - Recoverer
	mux.Use(middleware.Recoverer)
	//using middleware Nosurf - adds CSRF protection on all POST request
	mux.Use(NoSurf)
	//using middleware sessionLoad - this loads and saves the session on every request
	mux.Use(SessionLoad)

	//mux.Use(WriteToConsole)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals", handlers.Repo.Generals)
	mux.Get("/majors", handlers.Repo.Majors)

	//Availability routes
	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)

	mux.Get("/contact", handlers.Repo.Contact)

	//MakeReservation routes
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	//routing/putting the picture on your page by creating a file server
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

//Everytime someone hits the page we will write to the console
//Fmt.println with a message

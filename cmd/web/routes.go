package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/goldalee/golangprojects/GoWebDev/pkg/config"
	"github.com/goldalee/golangprojects/GoWebDev/pkg/handlers"
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
	return mux
}

//Everytime someone hits the page we will write to the console
//Fmt.println with a message

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/goldalee/golangprojects/bookings/internal/config"
	"github.com/goldalee/golangprojects/bookings/internal/handlers"
	"github.com/goldalee/golangprojects/bookings/internal/render"
)

const portNumber = ":8081"

var app config.AppConfig
var session *scs.SessionManager

// main is the main function
func main() {

	//change this to true when in production
	app.InProduction = false

	//Sessions
	session = scs.New()
	//we want it to last for 24 hours
	session.Lifetime = 24 * time.Hour
	//use cookies to store our sessions
	session.Cookie.Persist = true //we want the session to persist
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction //in production make sure it is true

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	//setting things up with our handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplate(&app)

	log.Printf("Starting application on port %v", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

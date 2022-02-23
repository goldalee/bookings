package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/goldalee/golangprojects/bookings/helpers"
	"github.com/goldalee/golangprojects/bookings/internal/config"
	"github.com/goldalee/golangprojects/bookings/internal/handlers"
	"github.com/goldalee/golangprojects/bookings/internal/models"
	"github.com/goldalee/golangprojects/bookings/internal/render"
)

const portNumber = ":8081"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {

	err := run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(fmt.Sprintf("Starting application on port #{portNumber}"))
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	//what am I going to put into the session

	//gob
	gob.Register(models.Reservation{})

	//change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

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
		log.Fatal("Cannot create template cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	//setting things up with our handlers
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplate(&app)
	helpers.NewHelpers(&app)
	return nil
}

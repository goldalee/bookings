package main

//Our own middleware that writes to the console
import (
	"net/http"

	"github.com/goldalee/golangprojects/bookings/helpers"
	"github.com/justinas/nosurf"
)

// func WriteToConsole(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("Hit the page")
// 		next.ServeHTTP(w, r)
// 	})

// }

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	//it uses cookies to make sure the token it generates is available on a per page basis
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction, //in production we will change this to true
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//Webservers by default are not state aware so we have to do this

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

//ne means not equal to - used on the about page in the if

func Auth(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if !helpers.IsAuthenticated(r){
			session.Put(r.Context(), "error", "Log in first!")
			http.Redirect(w,r,"/user/login", http.StatusSeeOther)
			return 
		}
		next.ServeHTTP(w,r)
	})
	}
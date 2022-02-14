package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/goldalee/golangprojects/bookings/internal/config"
	"github.com/goldalee/golangprojects/bookings/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	//gob
	gob.Register(models.Reservation{})
	//change this to true when in production
	testApp.InProduction = false

	//Sessions
	session = scs.New()
	//we want it to last for 24 hours
	session.Lifetime = 24 * time.Hour
	//use cookies to store our sessions
	session.Cookie.Persist = true //we want the session to persist
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction //in production make sure it is true

	testApp.Session = session

	app = &testApp
	os.Exit(m.Run())
}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

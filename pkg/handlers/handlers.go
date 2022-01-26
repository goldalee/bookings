package handlers

import (
	"net/http"

	"github.com/goldalee/golangprojects/bookings/pkg/config"
	"github.com/goldalee/golangprojects/bookings/pkg/models"
	"github.com/goldalee/golangprojects/bookings/pkg/render"
)

//variable that uses the repository type
var Repo *Repository

//the repository pattern - Repository is the reposity type
type Repository struct {
	App *config.AppConfig
}

//NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

//NewHandlers sets the reposity for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//performs some logic
	//putting some string unto about page
	stringMap := make(map[string]string)

	//information that will be passed
	stringMap["test"] = "Hello, again."

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	//sends data to the template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/goldalee/golangprojects/bookings/helpers"
	"github.com/goldalee/golangprojects/bookings/internal/config"
	"github.com/goldalee/golangprojects/bookings/internal/driver"
	"github.com/goldalee/golangprojects/bookings/internal/forms"
	"github.com/goldalee/golangprojects/bookings/internal/models"
	"github.com/goldalee/golangprojects/bookings/internal/render"
	"github.com/goldalee/golangprojects/bookings/internal/repository"
	"github.com/goldalee/golangprojects/bookings/internal/repository/dbrepo"
)

//variable that uses the repository type
var Repo *Repository

//the repository pattern - Repository is the reposity type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

//NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL, a),
		

	}
}

//NewHandlers sets the reposity for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	//m.DB.AllUsers()
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//sends data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Generals renders the room
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Reservations renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data, //empty reservation when the page is displayed for the first time
	})
}

// PostReservations handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd:=r.Form.Get("start_date")
	ed:=r.Form.Get("end_date")

	//2020-01--1 - 12/02 03:04:05PM '06 -0700
	
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil{
		helpers.ServerError(w, err)
	}

	endDate, err:= time.Parse(layout,ed)
	if err != nil{
		helpers.ServerError(w, err)
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil{
		helpers.ServerError(w, err)
	}




	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate: endDate,
		RoomID: roomID,
	}

	form := forms.New(r.PostForm)

	//form.Has("first_name", r)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	//write to database
	err= m.DB.InsertReservation(reservation)
	if err != nil{
		helpers.ServerError(w, err)
	}
	m.App.Session.Put(r.Context(), "reservation", reservation)
	//we do not want our user to enter the information twice
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// Post Availability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	//posting to the page
	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s", start, end)))
}

////////////////////////////
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

//AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//header
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

/////////////////////////////

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from summary")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation") //removing data from our session
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

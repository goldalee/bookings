package models

import "github.com/goldalee/bookings/internal/forms"

//only exist to be imported
//TemplateData holds data sent from handlers to templates- the sent data
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string //cross site request forgery token for security
	Flash     string //just posting success massage to user
	Warning   string //warning message
	Error     string //error message
	Form      *forms.Form
}

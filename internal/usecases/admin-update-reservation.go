package usecases

import (
	"log"
	"net/http"
	"strconv"

	"github.com/AthirsonSilva/golang-net-http-restapi/internal/forms"
	"github.com/AthirsonSilva/golang-net-http-restapi/internal/helpers"
	"github.com/AthirsonSilva/golang-net-http-restapi/internal/models"
)

func (r *Repository) AdminUpdateReservation(res http.ResponseWriter, req *http.Request) {
	log.Printf("[AdminUpdateReservation] Updating reservation...")
	var reservation models.Reservation

	err := req.ParseForm()
	if err != nil {
		log.Println(err)
		RedirectWithError(r, req, res, "Error parsing form", "/admin/reservations/all")
		return
	}

	id, err := strconv.Atoi(req.Form.Get("id"))
	if err != nil {
		log.Println(err)
		RedirectWithError(r, req, res, "Error parsing id", "/admin/reservations/all")
		return
	}

	startDate := req.Form.Get("start_date")
	endDate := req.Form.Get("end_date")

	var processed int
	strProcessed := req.Form.Get("processed")
	if strProcessed == "on" {
		processed = 1
	}

	form := forms.New(req.Form)
	form.Required("firstName", "lastName", "email", "phone", "start_date", "end_date")
	form.MinLength(2, req, "firstName", "lastName", "email", "phone", "start_date", "end_date")
	form.IsEmail("email")

	parsedStartDate := helpers.ConvertDateFromString(startDate, res)
	parsedEndDate := helpers.ConvertDateFromString(endDate, res)

	newReservation := models.Reservation{
		ID:        id,
		StartDate: parsedStartDate,
		EndDate:   parsedEndDate,
		Processed: processed,
	}

	err = r.Database.UpdateReservation(newReservation)
	if err != nil {
		log.Println(err)
		RedirectWithError(r, req, res, "Error updating reservation", "/admin/reservations/all")
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	http.Redirect(res, req, "/admin/reservations/all", http.StatusSeeOther)
}

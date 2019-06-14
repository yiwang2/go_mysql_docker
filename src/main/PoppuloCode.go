package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tickets", GetTickets).Methods("GET")
	router.HandleFunc("/ticket/{id}", GetSingleTicket).Methods("GET")
	router.HandleFunc("/ticket/{id}", ModifyTicket).Methods("PUT")
	router.HandleFunc("/ticket", CreateTicket).Methods("POST")
	router.HandleFunc("/ticket/{id}/status", ModifyTicketStatus).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func GetTickets(w http.ResponseWriter, r *http.Request) {

	results, err := GetAllTicketsFromDB()

	if err != nil {
		publishInterServerError(w, "Failed to get all tickets")
		return;
	}
	defer results.Close()
	tickets := CreateTicketIdsFromRows(results)
	json.NewEncoder(w).Encode(tickets)

}

func GetSingleTicket(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ticketId, result := validateInputId(params, w)
	if result {
		results, rowError := GetTicketInfoById(ticketId)
		if rowError != nil {
			publishInterServerError(w, "Failed to have ticket")
			return
		}
		defer results.Close()
		tickets := CreateTicketsFromRows(results)
		if tickets == nil || len(tickets) == 0 {
			publishInterServerError(w, "Failed to have ticket")
			return
		}

		json.NewEncoder(w).Encode(tickets[0])
	}
}

//after post create ticket, we are going to return an id of the new ticket
func CreateTicket(w http.ResponseWriter, r *http.Request) {
	var ticket Ticket
	err := json.NewDecoder(r.Body).Decode(&ticket)
	if validateCreateTicketInput(err, ticket, w) {
		ticket = ticket.CreateTicket(ticket.Status, ticket.Lines)
		saveErr := SaveTicket(ticket)
		if saveErr != nil {
			publishInterServerError(w, saveErr.Error())
			return
		}

		json.NewEncoder(w).Encode(Ticket{ticket.Id, nil, ""})
	}
}

func ModifyTicket(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	ticketId, result := validateInputId(params, w)
	if result {
		results, rowerr := GetTicketStatusById(ticketId)
		if rowerr != nil {
			publishInterServerError(w, "Failed to have ticket information")
			return
		}
		defer results.Close()
		tickets := CreateTicketIdAndStatusFromRow(results)
		if tickets == nil || len(tickets) == 0 {
			publishInterServerError(w, "Cannot find ticket with Id: "+string(ticketId))
			return
		}

		ticket := tickets[0]
		if ticket.Status == Checked {
			publishForbiddenError(w, "Cannot modify ticket due to its status")
			return
		}

		ticketLines, linesError := convertTicketLinesFromRequest(w, r)
		if !linesError {
			return
		}

		if validateTicketLines(ticketLines, w) {
			for i := 0; i < len(ticketLines); i++ {
				ticketLine := &ticketLines[i]
				ticketLine.LineResult = ticketLine.GetLineResult()
			}

			insertErr := AppendTicketLines(ticketId, ticketLines)
			if insertErr != nil {
				publishInterServerError(w, insertErr.Error())
				return
			}

			json.NewEncoder(w).Encode(Ticket{ticketId, nil, ""})
		}
	}
}

func ModifyTicketStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ticketId, result := validateInputId(params, w)
	if !result {
		return
	}

	ticket, result := convertTicketFromRequest(w, r)
	if !result {
		return
	}

	if !validateTicketStatusValue(ticket, w) {
		return
	}

	results, rowErr := GetTicketInfoById(ticketId)
	if rowErr != nil {
		publishInterServerError(w, "Failed to have ticket")
		return
	}
	defer results.Close()
	tickets := CreateTicketsFromRows(results)
	if tickets == nil || len(tickets) == 0 {
		publishInterServerError(w, "Failed to have ticket")
		return
	}

	if ticket.Status == tickets[0].Status {
		publishBadRequestError(w, "Tikcet is already at request status")
		return
	}

	updateError := UpDateTicketStatus(ticketId, ticket.Status)
	if updateError != nil {
		publishInterServerError(w, updateError.Error())
		return
	}

	json.NewEncoder(w).Encode(Ticket{ticketId, nil, ""})
}

func convertTicketLinesFromRequest(w http.ResponseWriter, r *http.Request) ([]TicketLine, bool) {
	var ticketLines []TicketLine
	linesErr := json.NewDecoder(r.Body).Decode(&ticketLines)

	if linesErr != nil {
		publishBadRequestError(w, linesErr.Error())
		return nil, false
	}

	return ticketLines, true
}

func convertTicketFromRequest(w http.ResponseWriter, r *http.Request) (Ticket, bool) {
	var ticket Ticket
	ticketErr := json.NewDecoder(r.Body).Decode(&ticket)

	if ticketErr != nil {
		publishBadRequestError(w, ticketErr.Error())
		return ticket, false
	}

	return ticket, true
}

func validateTicketLines(ticketLines []TicketLine, w http.ResponseWriter) bool {

	if ticketLines == nil || len(ticketLines) == 0 {
		publishBadRequestError(w, "Ticket lines are empty")
		return false
	}

	for i := 0; i < len(ticketLines); i++ {
		if !ticketLines[i].ValuesValid() {
			publishBadRequestError(w, "Bad ticket line values")
			return false
		}
	}

	return true
}

func validateTicketStatusValue(ticket Ticket, w http.ResponseWriter) bool {

	if ticket.Status != Checked && ticket.Status != Unchecked {
		publishBadRequestError(w, "Ticket status is invalid")
		return false
	}

	return true
}

func validateInputId(params map[string]string, w http.ResponseWriter) (int64, bool) {

	ticketId, err := strconv.ParseInt(string(params["id"]), 10, 64)

	if err != nil {
		publishBadRequestError(w, "Invalid Id: "+string(params["id"]))
		return -1, false
	}
	return ticketId, true
}

func validateCreateTicketInput(err error, ticket Ticket, w http.ResponseWriter) bool {
	if err != nil {
		publishBadRequestError(w, err.Error())
		return false
	}

	if ticket.Lines == nil || len(ticket.Lines) == 0 {
		publishBadRequestError(w, "Ticket lines are empty")
		return false
	}

	if ticket.Status != Checked && ticket.Status != Unchecked {
		publishBadRequestError(w, "Ticket status is invalid")
		return false
	}
	return true
}

func publishForbiddenError(w http.ResponseWriter, errorMessage string) {
	publisError(w, errorMessage, http.StatusForbidden)
}

func publishInterServerError(w http.ResponseWriter, errorMessage string) {
	publisError(w, errorMessage, http.StatusInternalServerError)
}

func publishBadRequestError(w http.ResponseWriter, errorMessage string) {
	publisError(w, errorMessage, http.StatusBadRequest)
}

func publisError(w http.ResponseWriter, errorMessage string, status int) {
	http.Error(w, errorMessage, status)
}

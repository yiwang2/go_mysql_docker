package main

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestRows(t *testing.T) {

	var ticketId int64 = 1;
	var ticketStatus string = string(Unchecked);

	t.Run("TestTicketRowCreation", func(t *testing.T) {

		var ticketLine *TicketLine = &TicketLine{[]TicketValue{2, 2, 2}, 5};
		lineStr, _ := ticketLine.ToString();
		var tckRow ticketRow = ticketRow{ticketId, ticketStatus, lineStr};
		if tckRow.id != ticketId {
			t.Error("Got: " + string(tckRow.id) + " Expected: " + string(ticketId))
		}
		if tckRow.status != ticketStatus {
			t.Error("Got: " + string(tckRow.status) + " Expected: " + ticketStatus)
		}
		if tckRow.line != lineStr {
			t.Error("Got: " + tckRow.line + " Expected: " + lineStr)
		}
	})

	//createTicketLineFromRow
	t.Run("TestCreateTicketLineFromRow", func(t *testing.T) {

		var ticketLine *TicketLine = &TicketLine{[]TicketValue{2, 2, 2}, 5};
		lineStr, _ := ticketLine.ToString();
		var tckRow ticketRow = ticketRow{0, "", lineStr};
		line, err :=createTicketLineFromRow(tckRow);
		if err != nil {
			t.Error("failed to createTicketLineFromRow:", err);
		}

		lineStrResult, _ := line.ToString()

		if  lineStrResult!= lineStr {
			t.Error("Got: " + lineStrResult + " Expected: " + lineStr)
		}
	})

	///createTicketFromRow
	t.Run("TestCreateTicketFromRow", func(t *testing.T) {

		var ticketLine *TicketLine = &TicketLine{[]TicketValue{2, 2, 2}, 5};
		lineStr, _ := ticketLine.ToString();
		var tckRow ticketRow = ticketRow{ticketId, ticketStatus, lineStr};
		tkt, err :=createTicketFromRow(tckRow);
		if err != nil {
			t.Error("failed to createTicketLineFromRow:", err);
		}

		lineStrResult, _ := tkt.Lines[0].ToString()

		if  lineStrResult!= lineStr {
			t.Error("Got: " + lineStrResult + " Expected: " + lineStr)
		}

		if  tkt.Id != ticketId{
			t.Error("Got: " + string(tkt.Id ) + " Expected: " + string(ticketId))
		}

		if  tkt.Status!= TicketStatus(ticketStatus) {
			t.Error("Got: " + string(tkt.Status) + " Expected: " + ticketStatus)
		}
	})
}

func TestCreateTicketIdAndStatusFromRow(t *testing.T) {

	var ticketId int64 = 1;
	var ticketStatus string = string(Unchecked);

	t.Run("TestCreateTicketIdAndStatusFromRow", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Error("failed to open sqlmock database:", err);
		}
		defer db.Close()
		rows := sqlmock.NewRows([]string{"id", "status"}).
			AddRow(ticketId, ticketStatus)
		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		rs, _ := db.Query("SELECT")
		defer rs.Close()

		if rs == nil  {
			t.Error("failed to query sqlmock database");
		}

		var tickets []Ticket = CreateTicketIdAndStatusFromRow(rs);

		if tickets == nil || len(tickets) != 1 {
			t.Error("failed to CreateTicketIdAndStatusFromRow");
		}

		var ticket Ticket = tickets[0];

		if ticket.Id != ticketId || ticket.Status != TicketStatus(ticketStatus) {
			t.Error("failed to create ticket based function CreateTicketIdAndStatusFromRow");
		}
	})

}


func TestCreateTicketFromRows(t *testing.T) {

	var ticketId int64 = 1;
	var ticketStatus string = string(Unchecked);
	var ticketLine *TicketLine = &TicketLine{[]TicketValue{2, 2, 2}, 5};
	lineStr, _ := ticketLine.ToString();

	t.Run("TestCreateTicketFromRows", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Error("failed to open sqlmock database:", err);
		}
		defer db.Close()
		rows := sqlmock.NewRows([]string{"id", "status", "line"}).
			AddRow(ticketId, ticketStatus, lineStr)
		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		rs, _ := db.Query("SELECT")
		defer rs.Close()

		if rs == nil  {
			t.Error("failed to query sqlmock database");
		}

		var tickets []Ticket = CreateTicketsFromRows(rs);

		if tickets == nil || len(tickets) != 1 {
			t.Error("failed to CreateTicketIdsFromRows");
		}

		var ticket Ticket = tickets[0];

		if ticket.Id != ticketId ||  ticket.Status != TicketStatus(ticketStatus) {
			t.Error("failed to create ticket based function CreateTicketsFromRows");
		}

		lineStrResult, _ := ticket.Lines[0].ToString();
		if lineStrResult != lineStr {
			t.Errorf("got %f want %f", lineStrResult, lineStr)
		}
	})
}

//CreateTicketsFromRows
func TestCreateTicketIdsFromRows(t *testing.T) {

	var ticketId int64 = 1;

	t.Run("TestCreateTicketIdsFromRows", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Error("failed to open sqlmock database:", err);
		}
		defer db.Close()
		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(ticketId)
		mock.ExpectQuery("SELECT").WillReturnRows(rows)

		rs, _ := db.Query("SELECT")
		defer rs.Close()

		if rs == nil  {
			t.Error("failed to query sqlmock database");
		}

		var tickets []Ticket = CreateTicketIdsFromRows(rs);

		if tickets == nil || len(tickets) != 1 {
			t.Error("failed to CreateTicketIdsFromRows");
		}

		var ticket Ticket = tickets[0];

		if ticket.Id != ticketId {
			t.Error("failed to create ticket based function CreateTicketIdsFromRows");
		}
	})
}

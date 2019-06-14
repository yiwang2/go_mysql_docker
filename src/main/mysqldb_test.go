package main

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestDB(t *testing.T) {

	var ticketId int64 = 1;
	var ticketStatus string = string(Unchecked);
	var ticketLine *TicketLine = &TicketLine{[]TicketValue{2, 2, 2}, 5};
	lineStr, _ := ticketLine.ToString();
	var ticket Ticket = Ticket{ticketId, []*TicketLine{ticketLine}, TicketStatus(ticketStatus)};

	t.Run("TestSaveTicket", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening database", err)
		}
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO tickets").WithArgs(ticketId, ticketStatus).WillReturnResult(sqlmock.NewResult(ticketId, 1))
		mock.ExpectExec("INSERT INTO ticket_line").WithArgs(ticketId, lineStr).WillReturnResult(sqlmock.NewResult(ticketId, 1))
		mock.ExpectCommit()

		if err = saveTicketIntoDB(db, ticket); err != nil {
			t.Errorf("error was not expected while updating stats: %s", err)
		}

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})

	t.Run("TestUpDateTicketStatus", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening database", err)
		}
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE tickets").WillReturnResult(sqlmock.NewResult(ticketId, 1))
		mock.ExpectCommit()

		if err = updateTctStatusInDB(db, ticketId, Unchecked); err != nil {
			t.Errorf("error was not expected while updating stats: %s", err)
		}

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})

	t.Run("TestAppendTicketLinesIntoDB", func(t *testing.T) {

		var ticketLineAppend TicketLine = TicketLine{[]TicketValue{2, 2, 2}, 5};

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening database", err)
		}
		defer db.Close()

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO ticket_line").WithArgs(ticketId, lineStr).WillReturnResult(sqlmock.NewResult(ticketId, 1))
		mock.ExpectCommit()

		if err = appendTicketLinesIntoDB(db, ticketId, []TicketLine{ticketLineAppend}); err != nil {
			t.Errorf("error was not expected while inserting lines: %s", err)
		}

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})

	t.Run("TestGetAllTicketsFromDB", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening database", err)
		}
		defer db.Close()
		sqlstaement:= "SELECT tickets.id From tickets"
		columns := []string{"id"}
		mock.ExpectBegin()
		mock.ExpectQuery(sqlstaement).WillReturnRows(sqlmock.NewRows(columns))

		_,err = getAllTicketRowsFromDB(db);

		if err != nil {
			t.Errorf("error was not expected while inserting lines: %s", err)
		}


		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})

	t.Run("TestGetTicketDetailsWithId", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening database", err)
		}
		defer db.Close()
		sqlstaement:= "SELECT tickets.id, tickets.status,ticket_line.line FROM tickets INNER JOIN ticket_line ON tickets.id = ticket_line.ticket_id and tickets.id=?"
		columns := []string{"id", "status", "line"}
		mock.ExpectBegin()
		mock.ExpectQuery(sqlstaement).WithArgs(ticketId).WillReturnRows(sqlmock.NewRows(columns))

		_,err = getTicketDetailsWithId(db, ticketId);

		if err != nil {
			t.Errorf("error was not expected while inserting lines: %s", err)
		}


		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})

	t.Run("TestGetTicketStatusWithId", func(t *testing.T) {

		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening database", err)
		}
		defer db.Close()
		sqlstaement:= "SELECT tickets.id, tickets.status FROM tickets where tickets.id=?"
		columns := []string{"id", "status"}
		mock.ExpectBegin()
		mock.ExpectQuery(sqlstaement).WithArgs(ticketId).WillReturnRows(sqlmock.NewRows(columns))

		_,err = getTicketStatus(db, ticketId);

		if err != nil {
			t.Errorf("error was not expected while inserting lines: %s", err)
		}


		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})

}

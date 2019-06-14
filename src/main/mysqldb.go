package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

const drivername = "mysql";
const datasourcename = "testuser:testuser@tcp(db:3306)/poppulodb";


func AppendTicketLines (id int64, lines []TicketLine) error{
	db, err := sql.Open(drivername, datasourcename)
	defer db.Close();
	if err != nil {
		return err;
	}

	return appendTicketLinesIntoDB(db, id, lines);
}

func appendTicketLinesIntoDB (db *sql.DB,id int64, lines []TicketLine) error{
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	var insertStatement string = "INSERT INTO ticket_line (ticket_id, line) values ";
	vals := []interface{}{}
	for _, line :=range lines {
		var value, tostrerr = line.ToString();
		if tostrerr != nil {
			continue;
		}

		insertStatement += "(?, ?),"
		vals = append(vals, id, value);
	}
	//remove last ','
	insertStatement = insertStatement[0:len(insertStatement)-1]
	if _, err = tx.Exec(insertStatement, vals...); err != nil {
		return err
	}

	return nil;
}

func GetTicketStatusById (id int64) (*sql.Rows, error){
	db, err := sql.Open(drivername, datasourcename)
	defer db.Close();
	if err != nil {
		return nil, err;
	}

	return getTicketStatus(db, id);
}

func getTicketStatus (db *sql.DB, id int64) (*sql.Rows, error){
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		switch err {
		case nil:
		default:
			tx.Rollback()
		}
	}()

	sqlstmt := "SELECT tickets.id, tickets.status FROM tickets where tickets.id=?"

	return tx.Query(sqlstmt, id)
}


func GetTicketInfoById(id int64) (*sql.Rows, error) {
	db, err := sql.Open(drivername, datasourcename)
	defer db.Close();
	if err != nil {
		return nil, err;
	}

	return getTicketDetailsWithId(db, id);
}

func getTicketDetailsWithId (db *sql.DB, id int64) (*sql.Rows, error){
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		switch err {
		case nil:
		default:
			tx.Rollback()
		}
	}()

	sqlstmt := "SELECT tickets.id, tickets.status,ticket_line.line FROM tickets INNER JOIN ticket_line ON tickets.id = ticket_line.ticket_id and tickets.id=?"

	return tx.Query(sqlstmt, id)
}

func GetAllTicketsFromDB() (*sql.Rows, error) {
	db, err := sql.Open(drivername, datasourcename)
	defer db.Close();
	if err != nil {
		return nil, err;
	}

	return getAllTicketRowsFromDB(db);
}

func getAllTicketRowsFromDB (db *sql.DB) (*sql.Rows, error){
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		switch err {
		case nil:
		default:
			tx.Rollback()
		}
	}()

	return tx.Query("SELECT tickets.id From tickets")
}

func UpDateTicketStatus (id int64, status TicketStatus) error {
	db, err := sql.Open(drivername, datasourcename)
	defer db.Close();
	if err != nil {
		return err;
	}

	return updateTctStatusInDB(db, id, status)
}

func updateTctStatusInDB (db *sql.DB, id int64, status TicketStatus) error{
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if _, err = tx.Exec("UPDATE tickets SET status=? where id=?", status, id); err != nil {
		return err
	}

	return nil
}

func SaveTicket(ticket Ticket) error {
	db, err := sql.Open(drivername, datasourcename)
	defer db.Close();
	if err != nil {
		return err;
	}

	return saveTicketIntoDB(db, ticket)
}

func saveTicketIntoDB (db *sql.DB,ticket Ticket) error{
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	if _, err = tx.Exec("INSERT INTO tickets VALUES( ?, ? )", ticket.Id, ticket.Status); err != nil {
		return err
	}
	for _, ticketLine := range ticket.Lines {
		lineResult, lineErr := ticketLine.ToString()
		if (lineErr != nil) {
			return lineErr;
		}
		if _, err = tx.Exec("INSERT INTO ticket_line VALUES( ?, ? )", ticket.Id, lineResult); err != nil {
			return err
		}
	}

	return nil;
}

package main

import (
	"database/sql"
	"encoding/json"
	"sort"
	"log"
)

type ticketRow struct {
	id     int64
	status string
	line   string
}

func CreateTicketIdAndStatusFromRow(rows *sql.Rows) []Ticket {
	var tickets []Ticket
	for rows.Next() {
		var row ticketRow
		err := rows.Scan(&row.id, &row.status)
		if err != nil {
			continue
		}
		tickets = append(tickets, Ticket{row.id, nil, TicketStatus(row.status)})
	}

	return tickets
}


//get all ids only when query everything
func CreateTicketIdsFromRows(rows *sql.Rows) []Ticket {
	var tickets []Ticket
	for rows.Next() {
		var row ticketRow
		err := rows.Scan(&row.id)
		if err != nil {
			continue
		}
		tickets = append(tickets, Ticket{row.id, nil, ""})
	}

	return tickets
}

//manage rows back from database - used for get single ticket details
func CreateTicketsFromRows(rows *sql.Rows) []Ticket {

	var tickets []Ticket
	ticketsMap := make(map[int64]*Ticket)
	for rows.Next() {
		var row ticketRow
		err := rows.Scan(&row.id, &row.status, &row.line)
		if err != nil {
			continue
		}

		if ticket, ok := ticketsMap[row.id]; ok {
			line, lineErr := createTicketLineFromRow(row)
			if lineErr != nil {
				log.Fatal(lineErr.Error())
				continue
			}
			ticket.AppendLines([]*TicketLine {&line})
		} else {
			newTicket, ticketErr := createTicketFromRow(row)
			if ticketErr != nil {
				log.Fatal(ticketErr.Error())
				continue
			}
			ticketsMap[newTicket.Id] = &newTicket
		}
	}

	for  _, ticket := range ticketsMap {
		sort.Sort(SortByResult{ticket.Lines})
		tickets = append(tickets, *ticket)
	}

	return tickets
}

func createTicketLineFromRow (row ticketRow) (TicketLine, error){
	var line TicketLine
	err := json.Unmarshal([]byte(row.line), &line)

	if err != nil {
		return TicketLine{},err
	}

	return line, nil
}

func createTicketFromRow (row ticketRow) (Ticket, error){
	line, err := createTicketLineFromRow(row)

	if err != nil {
		return Ticket{},err
	}

	return  Ticket{row.id,TicketLines{&line}, TicketStatus(row.status)}, nil
}

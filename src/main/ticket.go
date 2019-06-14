package main

import (
	"errors"
	"strconv"
	"time"
	"sort"
	"encoding/json"
)

//each of which has a value of 0, 1, or 2.
type TicketValue int

const (
	ValueZero TicketValue = iota // value --> 0
	ValueOne                     // value --> 1
	ValueTwo                     // value --> 2
)

//we will need to have the ability to check the status of lines on a ticket
type TicketStatus string

const (
	Checked   TicketStatus = "CHECKED"
	Unchecked TicketStatus = "UNCHECKED"
)

type TicketLines []*TicketLine

//a series of lines on a ticket
type Ticket struct {
	Id     int64        `json:"id,omitempty"`
	Lines  TicketLines  `json:"lines,omitempty"`
	Status TicketStatus `json:"status,omitempty"`
}

//create a ticket
func (ticket *Ticket) CreateTicket(status TicketStatus, lines TicketLines) Ticket {
	//assuming id is a time stamp
	ticketId := time.Now().UnixNano() / int64(time.Millisecond)

	for _, line := range lines {
		line.LineResult = line.GetLineResult()
	}
	sort.Sort(SortByResult{lines})

	return Ticket{ticketId, lines, status}
}

//Once the status of a ticket has been checked it should not be possible to amend.
func (ticket *Ticket) IsAmendable() bool {
	return ticket.Status == Unchecked
}

//this is only for creation or after query
func (ticket *Ticket) AppendLines(lines TicketLines) {
	ticket.Lines = append(ticket.Lines, lines...)
}

//It should be possible for a ticket to be amended with n additional lines.
//return num of lines or error if there is something wrong
func (ticket *Ticket) AddLines(lines TicketLines) (int, error) {

	if ticket.IsAmendable() {
		ticket.Lines = append(ticket.Lines, lines...)
		return len(ticket.Lines), nil
	} else {
		return -1, errors.New("Ticket with id: " + strconv.FormatInt(ticket.Id, 10) + " is not ammendable")
	}
}

//We would like the lines sorted into outcomes before being returned.
//usage sort.Sort(SortByResult{ticket.Lines})
func (t TicketLines) Len() int      { return len(t) }
func (t TicketLines) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

//prefer descending
type SortByResult struct{ TicketLines }

func (p SortByResult) Less(i, j int) bool {
	return p.TicketLines[i].GetLineResult() > p.TicketLines[j].GetLineResult()
}

type TicketLine struct {
	Values     []TicketValue `json:"values"`
	LineResult int           `json:"result"`
}

func (line *TicketLine) ValuesValid() bool {
	if line.Values == nil || len(line.Values) == 0 {
		return false
	}

	for _, value := range line.Values {
		if value > ValueTwo || value < ValueZero {
			return false
		}
	}

	return true
}

func (line *TicketLine) ToString() (string, error) {
	b, err := json.Marshal(line)
	if err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

func (line *TicketLine) GetLineResult() int {
	lineSum := int(line.Values[0])
	firstValue:= int(line.Values[0])
	allSame:= true
	firstDifferent := true
	for i := 1; i < len(line.Values); i++ {
		lineSum = lineSum + int(line.Values[i])
		if int(line.Values[i]) != firstValue {
			allSame = false
		} else {
			//we want both 2nd and 3rd numbers are different from 1st one
			//if code jump into this block, we have 1 num is equal to first
			firstDifferent = false
		}
	}
	//if the sum of the values on a line is 2, the result for that line is 10
	if lineSum == 2 {
		return 10
	} else if allSame { //if they are all the same, the result is 5
		return 5
	} else if firstDifferent { //both 2nd and 3rd numbers are different from the 1st=
		return 1
	} else {
		return 0
	}
}

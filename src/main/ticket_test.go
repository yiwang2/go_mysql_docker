package main

import (
	"testing"
)

func TestTicket (t *testing.T) {
	var ticketLine *TicketLine = &TicketLine{[]TicketValue{2,2,2},5};
	var ticket *Ticket = &Ticket{};


	t.Run("TestCreateTicket", func(t *testing.T) {
		var tkt Ticket = ticket.CreateTicket("UNCHECKED", []*TicketLine {ticketLine})

		if tkt.Status != Unchecked {
			t.Error("Got: "+tkt.Status+" Expected: "+Unchecked)
		}

		if len(tkt.Lines) != 1 {
			t.Error("Number of ticket line is wrong")
		}

		if !tkt.IsAmendable() {
			t.Error("Ticket status as UNCHECKED should be allowed ammended")
		}
	})

	t.Run("TestTicketStatus", func(t *testing.T) {
		var tkt Ticket = ticket.CreateTicket("CHECKED", []*TicketLine {ticketLine})
		if tkt.IsAmendable() {
			t.Error("Ticket status as CHECKED should not be allowed ammended")
		}
	})

	t.Run("TestTicketAddLines", func(t *testing.T) {
		var tkt Ticket = ticket.CreateTicket("CHECKED", []*TicketLine {ticketLine})
		var anotherticketLine *TicketLine = &TicketLine{[]TicketValue{2,2,2},5};

		_, err := tkt.AddLines([]*TicketLine {anotherticketLine})
		if err == nil {
			t.Error("Ticket status as CHECKED should not be allowed adding lines")
		}

		var tp *Ticket = &tkt;
		tp.Status = Unchecked;
		n, err :=tkt.AddLines([]*TicketLine {anotherticketLine})

		if n != len(tkt.Lines) {
			t.Error("Failed to add ticket lines")
		}

	})

	t.Run("TestTicketAppendLines", func(t *testing.T) {
		var tkt Ticket = ticket.CreateTicket("CHECKED", []*TicketLine {ticketLine})
		var anotherticketLine *TicketLine = &TicketLine{[]TicketValue{2,2,2},5};
		tkt.AppendLines([]*TicketLine {anotherticketLine})

		if len(tkt.Lines) != 2 {
			t.Error("Failed to append ticket lines")
		}

	})

}


func TestTicketLine(t *testing.T) {

	var ticketLine *TicketLine = &TicketLine{[]TicketValue{2,2,2},5};

	var ticketLineDummy *TicketLine = &TicketLine{[]TicketValue{3,3,3},5};

	t.Run("TestTicketLineDummy", func(t *testing.T) {
		r1 := ticketLine.ValuesValid();
		r2 :=ticketLineDummy.ValuesValid();

		if r1 == false {
			t.Error("ticket validation failed")
		}

		if r2 == true {
			t.Error("ticket validation failed")
		}
	})

	t.Run("TestTicketLineToString", func(t *testing.T) {
		got, err:= ticketLine.ToString()
		if err != nil {
			t.Error("Error is not nil")
		}

		var expected string = `{"values":[2,2,2],"result":5}`;

		if got != expected {
			t.Error("Got: "+got+" Expected: "+expected)
		}
	})

	t.Run("TestTicketGetResult", func(t *testing.T) {

		// expected result as 10 because sum as 2
		ticketLine.Values = []TicketValue{0,2,0}
		ticketLine.LineResult = 10;
		got1 := ticketLine.GetLineResult();
		if got1 != ticketLine.LineResult {
			t.Errorf("got %.2f want %.2f", got1, ticketLine.LineResult)
		}

		// expected result as 5 because all are same
		ticketLine.Values = []TicketValue{2,2,2}
		ticketLine.LineResult = 5;
		got2 := ticketLine.GetLineResult();
		if got2 != ticketLine.LineResult {
			t.Errorf("got %d want %d", got2, ticketLine.LineResult)
		}

		// expected result as 0
		ticketLine.Values = []TicketValue{2,2,1}
		ticketLine.LineResult = 0;
		got3 := ticketLine.GetLineResult();
		if got3 != ticketLine.LineResult {
			t.Errorf("got %d want %d", got3, ticketLine.LineResult)
		}

		// expected result as 0
		ticketLine.Values = []TicketValue{0,2,1}
		ticketLine.LineResult = 1;
		got4 := ticketLine.GetLineResult();
		if got4 != ticketLine.LineResult {
			t.Errorf("got %d want %d", got4, ticketLine.LineResult)
		}
	})

}
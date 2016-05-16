package relaymail

import (
	"testing"
)

func TestHeader(t *testing.T) {
	m := NewInMail()
	m.MsgId = "msgid1234"
	m.EnvelopeFrom = "Mail From: <from@example.com>"
	m.DataHeader = append(m.DataHeader, "From: from")
	m.DataHeader = append(m.DataHeader, "To: to")
	m.DataHeader = append(m.DataHeader, "Date: date")
	m.DataHeader = append(m.DataHeader, "Message-id: id")
	m.DataHeader = append(m.DataHeader, "Subject: test")
	m.DataHeader = append(m.DataHeader, "")

	m.AddAutoHeader()

	if len(m.DataHeader) != 6 {
		t.Error("add auto header fail", len(m.DataHeader))
	}
	if m.DataHeader[0] != "From: from" {
		t.Error("not subject[%s]", m.DataHeader[0])
	}
	if m.DataHeader[1] != "To: to" {
		t.Error("not from[%s]", m.DataHeader[1])
	}
	if m.DataHeader[2][0:6] != "Date: " {
		t.Error("not from[%s]", m.DataHeader[2])
	}
	if m.DataHeader[3] != "Message-id: id" {
		t.Error("not message-id[%s]", m.DataHeader[3])
	}
	if len(m.DataHeader[5]) != 0 {
		t.Error("not last header[%s]", m.DataHeader[4])
	}
}
func TestAddAutoHeader(t *testing.T) {
	m := NewInMail()
	m.MsgId = "msgid1234"
	m.EnvelopeFrom = "Mail From: <from@example.com>"
	m.DataHeader = append(m.DataHeader, "Subject: test")
	m.DataHeader = append(m.DataHeader, "")

	m.AddAutoHeader()

	if len(m.DataHeader) != 5 {
		t.Error("add auto header fail", len(m.DataHeader))
	}
	if m.DataHeader[0] != "Subject: test" {
		t.Error("not subject[%s]", m.DataHeader[0])
	}
	if m.DataHeader[1] != "From: <from@example.com>" {
		t.Error("not from[%s]", m.DataHeader[1])
	}
	if m.DataHeader[2][0:6] != "Date: " {
		t.Error("not from[%s]", m.DataHeader[2])
	}
	if m.DataHeader[3] != "Message-ID: <msgid1234>" {
		t.Error("not message-id[%s]", m.DataHeader[3])
	}
	if len(m.DataHeader[4]) != 0 {
		t.Error("not last header[%s]", m.DataHeader[4])
	}
}

package relaymail

import (
	"testing"
)

func TestAddAutoHeader(t *testing.T) {
	m := NewInMail()
	m.EnvelopeFrom = "Mail From: <from@example.com>"
	m.EnvelopeTo = append(m.EnvelopeTo, "Rcpt To: <to@example.com>")
	m.DataHeader = append(m.DataHeader, "Subject: test")

	m.AddAutoHeader()
}

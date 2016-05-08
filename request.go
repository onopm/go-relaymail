package relaymail

import ()

type InMail struct {
	EnvelopeFrom string
	EnvelopeTo   []string
	Data         []string
}

func NewInMail() *InMail {
	return &InMail{}
}

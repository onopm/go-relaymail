package relaymail

import (
	"fmt"
	"strings"
	"time"
)

type InMail struct {
	EnvelopeFrom string
	EnvelopeTo   []string
	DataHeader   []string
	DataBody     []string
}

func NewInMail() *InMail {
	return &InMail{}
}

var date_layout = "Mon, 2 Jan 2006 15:04:05 +0900 (JST)"

func (m *InMail) AddAutoHeader() {
	exists_from := false
	exists_date := false

	for i := 0; i < len(m.DataHeader); i++ {

		switch {
		case len(m.DataHeader[i]) > 4:
			h := strings.ToUpper(m.DataHeader[i][0:5])
			if h == "FROM:" {
				exists_from = true
			} else if h == "DATE:" {
				exists_date = true
			}

		case len(m.DataHeader[i]) < 1:
			lastHeader := m.DataHeader[i]
			m.DataHeader = m.DataHeader[0:(i - 1)]

			if exists_from == false {
				fromHeader := fmt.Sprintf("From: %s", m.EnvelopeFrom)
				m.DataHeader = append(m.DataHeader, fromHeader)
				warnf("add header|%s", fromHeader)
			}
			if exists_date == false {
				now := time.Now()
				dateHeader := fmt.Sprintf("Date: %s", now.Format(date_layout))
				m.DataHeader = append(m.DataHeader, dateHeader)
				warnf("add header|%s", dateHeader)
			}
			m.DataHeader = append(m.DataHeader, lastHeader)
			return
		}
	}
}

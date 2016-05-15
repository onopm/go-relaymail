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
	MsgId        string
}

func NewInMail() *InMail {
	return &InMail{}
}

var date_layout = "Mon, 2 Jan 2006 15:04:05 +0900 (JST)"

func (m *InMail) AddAutoHeader() {
	exist_from := false
	exist_date := false
	exist_msgid := false

	for i := 0; i < len(m.DataHeader); i++ {

		switch {
		case len(m.DataHeader[i]) > 4:
			h := strings.ToUpper(m.DataHeader[i][0:5])
			if h == "FROM:" {
				exist_from = true
			} else if h == "DATE:" {
				exist_date = true
			} else if h == "MESSG" {
				//TODO:
				if len(m.DataHeader[i]) > 10 {
					h := strings.ToUpper(m.DataHeader[i][0:11])
					if h == "MESSAGE-ID:" {
						exist_msgid = true
					}
				}
			}
		case len(m.DataHeader[i]) < 1:
			lastHeader := m.DataHeader[i]
			m.DataHeader = m.DataHeader[0:(i - 1)]

			if exist_from == false {
				eFrom := strings.SplitN(m.EnvelopeFrom, ":", 2)
				addHeader := fmt.Sprintf("From: %s", eFrom[1])
				m.DataHeader = append(m.DataHeader, addHeader)
				warnf("add header|%s", addHeader)
			}
			if exist_date == false {
				now := time.Now()
				addHeader := fmt.Sprintf("Date: %s", now.Format(date_layout))
				m.DataHeader = append(m.DataHeader, addHeader)
				warnf("add header|%s", addHeader)
			}
			if exist_msgid == false {
				addHeader := fmt.Sprintf("Message-ID: <%s>", m.MsgId)
				m.DataHeader = append(m.DataHeader, addHeader)
				warnf("add header|%s", addHeader)
			}
			m.DataHeader = append(m.DataHeader, lastHeader)
			return
		}
	}
}

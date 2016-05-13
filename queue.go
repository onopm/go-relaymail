package relaymail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/smtp"
	"strings"
)

type Queue struct {
	In      chan *InMail
	Stop    chan bool
	NextMTA string
}

func newQueue() *Queue {
	q := &Queue{}
	q.In = make(chan *InMail)
	q.Stop = make(chan bool)

	return q
}

func (q *Queue) serv() {
	infof("queue start")

	go func() {
	loop:
		for {
			select {
			case m := <-q.In:
				infof("enqueue [%s]", m.EnvelopeFrom)
				//TODO: goroutine or connection keep
				q.send(m)
			case <-q.Stop:
				infof("queue stop")
				break loop
			}
		}
	}()
}

func (q *Queue) saveJson(m *InMail) {

	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	//b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	//b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
	//b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	fmt.Println(string(b))

	//TODO: save file ?
}

var date_layout = "Mon,  2 Jan 2006 15:04:05 +0900 (JST)"

func (q *Queue) send(m *InMail) {

	infof("send start")

	c, err := smtp.Dial(q.NextMTA)
	if err != nil {
		warnf("%s", err)
		return
	}
	defer c.Close()

	sFrom := strings.Index(m.EnvelopeFrom, ":")
	if sFrom > 0 {
		c.Mail(m.EnvelopeFrom[sFrom:])
	}
	for i := 0; i < len(m.EnvelopeTo); i++ {
		sTo := strings.Index(m.EnvelopeTo[i], ":")
		if sTo > 0 {
			c.Rcpt(m.EnvelopeTo[i][sTo:])
		}
	}
	wc, err := c.Data()
	if err != nil {
		warnf("%s", err)
		return
	}
	defer wc.Close()

	//TODO:
	buf := bytes.NewBufferString(m.Data[0] + "\r\n")
	for i := 1; i < len(m.Data); i++ {
		buf.WriteString(m.Data[i] + "\r\n")
	}
	if _, err = buf.WriteTo(wc); err != nil {
		warnf("%s", err)
	}
}

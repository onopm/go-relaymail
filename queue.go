package relaymail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type Queue struct {
	In      chan *InMail
	Stop    chan bool
	NextMTA string
	Dir     string
}

func newQueue() *Queue {
	q := &Queue{}
	q.In = make(chan *InMail)
	q.Stop = make(chan bool)

	return q
}

func (q *Queue) serv() {
	infof("queue start. dir[%s]", q.Dir)
	d, err := os.Stat(q.Dir)
	if err != nil {
		warnf("queue dir fail[%s]", err)
		infof("create dir[%s]", q.Dir)
		os.Mkdir(q.Dir, 0700)
		d, _ = os.Stat(q.Dir)
	}
	if !d.IsDir() {
		warnf("queue dir not dir")
		os.Exit(1)
	}

	go func() {
	loop:
		for {
			select {
			case m := <-q.In:
				infof("enqueue [%s]", m.EnvelopeFrom)
				//TODO: goroutine or connection keep
				m.AddAutoHeader()
				err := q.send(m)
				if err != nil {
					warnf("send fail: %s", err)
				}
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

func (q *Queue) send(m *InMail) error {

	if len(q.NextMTA) < 1 {
		infof("send drop")
		return nil
	}
	infof("send start")

	c, err := smtp.Dial(q.NextMTA)
	if err != nil {
		warnf("%s", err)
		return fmt.Errorf("%s", err)
	}
	defer c.Close()

	sFrom := strings.Index(m.EnvelopeFrom, ":")
	if sFrom > 0 {
		err := c.Mail(m.EnvelopeFrom[sFrom:])
		if err != nil {
			warnf("%s", err)
			return fmt.Errorf("%s", err)
		}
	}
	for i := 0; i < len(m.EnvelopeTo); i++ {
		sTo := strings.Index(m.EnvelopeTo[i], ":")
		if sTo > 0 {
			err := c.Rcpt(m.EnvelopeTo[i][sTo:])
			if err != nil {
				warnf("%s", err)
				return fmt.Errorf("%s", err)
			}
		}
	}
	wc, err := c.Data()
	if err != nil {
		warnf("%s", err)
		return fmt.Errorf("%s", err)
	}
	defer wc.Close()

	//TODO:
	buf := bytes.NewBufferString(m.DataHeader[0] + "\r\n")
	for i := 1; i < len(m.DataHeader); i++ {
		buf.WriteString(m.DataHeader[i] + "\r\n")
	}
	for i := 0; i < len(m.DataBody); i++ {
		buf.WriteString(m.DataBody[i] + "\r\n")
	}
	if _, err = buf.WriteTo(wc); err != nil {
		warnf("%s", err)
		return fmt.Errorf("%s", err)
	}
	return nil
}

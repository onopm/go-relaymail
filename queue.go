package relaymail

import (
	"encoding/json"
	"fmt"
)

type Queue struct {
	In   chan *InMail
	Stop chan bool
}

func newQueue() *Queue {
	q := &Queue{}
	q.In = make(chan *InMail)
	q.Stop = make(chan bool)

	return q
}

func (q *Queue) serv() {
	fmt.Println("queue start")

	go func() {
	loop:
		for {
			select {
			case m := <-q.In:
				fmt.Printf("enqueue [%s]\n", m.EnvelopeFrom)
				break loop
			case <-q.Stop:
				fmt.Println("queue stop")
				return
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

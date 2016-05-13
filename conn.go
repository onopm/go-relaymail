package relaymail

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

type conn struct {
	server *Server  //set from Server
	rwc    net.Conn //set from Server

	id         string
	remoteAddr string
	//tlsState   *tls.ConnectionState
	//werr error
	//r          *connReader
	//bufr       *bufio.Reader
	//bufw       *bufio.Writer
	tp         *textproto.Reader
	lastMethod string
	isData     bool
	//mu         sync.Mutex
	//hijackedv bool
	mail []*InMail
}

func (c *conn) close() {
	c.rwc.Close()
}

func (c *conn) readData(q *Queue) error {
	lines, err := c.tp.ReadDotLines()
	if err != nil {
		warnf("%s", err)
		return err
	}

	m := c.mail[len(c.mail)-1]
	headerPart := true
	for i := 0; i < len(lines); i++ {
		infof("[%s] recv[%s]", c.id, lines[i])
		if headerPart == true {
			m.DataHeader = append(m.DataHeader, lines[i])
		} else {
			m.DataBody = append(m.DataBody, lines[i])
		}
		if len(lines[i]) < 1 {
			headerPart = false
		}
	}

	//TODO: save sync
	q.In <- m

	c.lastMethod = "DOT"
	c.isData = false

	return nil
}

func (c *conn) readRequest() error {
	line, err := c.tp.ReadLine()
	if err != nil {
		warnf("%s", err)
		return err
	}
	infof("[%s] recv[%s]", c.id, line)

	s1 := strings.Index(line, " ")
	if s1 < 0 {
		c.lastMethod = strings.ToUpper(line)
	} else {
		c.lastMethod = strings.ToUpper(line[:s1])
	}

	switch c.lastMethod {
	case "MAIL":
		m := NewInMail()
		m.EnvelopeFrom = line
		c.mail = append(c.mail, m)
	case "RCPT":
		m := c.mail[len(c.mail)-1]
		m.EnvelopeTo = append(m.EnvelopeTo, line)
	case "DATA":
		c.isData = true
	default:
		c.isData = false
	}

	return nil
}

func (c *conn) serve(q *Queue) {
	c.remoteAddr = c.rwc.RemoteAddr().String()
	c.isData = false
	c.id = createId()
	infof("[%s] connection from %v", c.id, c.remoteAddr)
	defer func() {
		infof("[%s] close connection", c.id)
		c.close()
	}()

	fmt.Fprintf(c.rwc, "220 greeting\r\n")

	c.tp = textproto.NewReader(bufio.NewReader(c.rwc))

	for {
		var err error
		if c.isData == false {
			err = c.readRequest()
		} else {
			err = c.readData(q)
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		switch c.lastMethod {
		case "HELO":
			fmt.Fprintf(c.rwc, "250 helo ok\r\n")
		case "EHLO":
			fmt.Fprintf(c.rwc, "250 ehlo ok\r\n")
		case "MAIL":
			fmt.Fprintf(c.rwc, "250 mail ok\r\n")
		case "RCPT":
			fmt.Fprintf(c.rwc, "250 rcpt ok\r\n")
		case "DATA":
			fmt.Fprintf(c.rwc, "354 data ok\r\n")
		case "DOT":
			fmt.Fprintf(c.rwc, "250 data ok\r\n")
		case "QUIT":
			fmt.Fprintf(c.rwc, "221 OK\r\n")
			return
		default:
			fmt.Fprintf(c.rwc, "500 unknown\r\n")
		}
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	layoutMilli = "20060102T150405.000"
	layoutMicro = "20060102T150405.000000"
	layoutNano  = "20060102T150405.000000000"
)

func createId() string {
	//ZZZ  46655
	//100   1296
	var x = rand.Intn(45360) + 1296
	strconv.FormatInt(int64(x), 36)

	//TODO:
	t := time.Now()
	id := t.Format(layoutMicro) + strings.ToUpper(strconv.FormatInt(int64(x), 36))

	return id
}

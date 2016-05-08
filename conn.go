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
}

func (c *conn) close() {
	c.rwc.Close()
}

func (c *conn) readData() error {
	lines, err := c.tp.ReadDotLines()
	if err != nil {
		fmt.Println(err)
		return err
	}

	for i := 0; i < len(lines); i++ {
		fmt.Printf("[%s] recv[%s]\n", c.id, lines[i])
	}

	c.lastMethod = "DOT"
	c.isData = false
	return nil
}

func (c *conn) readRequest() error {
	line, err := c.tp.ReadLine()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("[%s] recv[%s]\n", c.id, line)

	s1 := strings.Index(line, " ")
	if s1 < 0 {
		c.lastMethod = strings.ToUpper(line)
	} else {
		c.lastMethod = strings.ToUpper(line[:s1])
	}

	switch c.lastMethod {
	case "DATA":
		c.isData = true
	default:
		c.isData = false
	}

	return nil
}

func (c *conn) serve() {
	c.remoteAddr = c.rwc.RemoteAddr().String()
	c.isData = false
	c.id = createId()
	fmt.Printf("[%s] connection from %v\n", c.id, c.remoteAddr)
	defer func() {
		fmt.Printf("[%s] close connection\n", c.id)
		c.close()
	}()

	fmt.Fprintf(c.rwc, "220 greeting\r\n")

	c.tp = textproto.NewReader(bufio.NewReader(c.rwc))

	for {
		if c.isData == false {
			err := c.readRequest()
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			err := c.readData()
			if err != nil {
				fmt.Println(err)
				return
			}
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

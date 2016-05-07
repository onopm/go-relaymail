package relaymail

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"strings"
)

type conn struct {
	server *Server  //set from Server
	rwc    net.Conn //set from Server

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
	fmt.Printf("recv DATA[%v]\n", lines)

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
	fmt.Printf("recv[%s]\n", line)

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
	fmt.Printf("connection from %v\n", c.remoteAddr)
	defer func() {
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

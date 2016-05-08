package relaymail

import (
	"fmt"
	"net"
	"time"
)

type Server struct {
	Addr string
	//Handler        Handler
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	//TLSConfig      *tls.Config
	//TLSNextProto map[string]func(*Server, *tls.Conn, Handler)
	//ConnState func(net.Conn, ConnState)
	//ErrorLog *log.Logger
	//nextProtoOnce     sync.Once
	//nextProtoErr      error
}

func (srv *Server) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: srv,
		rwc:    rwc,
	}
	return c
}

func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":smtp"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	fmt.Println("listen: " + addr)

	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure

	q := newQueue()
	q.serv()

	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				fmt.Printf("smtp: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0
		c := srv.newConn(rw)
		go c.serve(q)
	}
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func ListenAndServe(addr string) error {
	server := &Server{Addr: addr}
	return server.ListenAndServe()
}

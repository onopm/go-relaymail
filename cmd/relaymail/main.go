package main

import (
	"flag"
	"fmt"
	"os"

	relaymail "github.com/onopm/go-relaymail"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: relaymail next-mta\n")
		flag.PrintDefaults()
	}
}

func main() {
	var (
		host   string
		port   int
		listen string
	)
	//TODO: ex) use "-listen :10025"
	flag.StringVar(&host, "host", "", "listen IP Addr.")
	flag.IntVar(&port, "port", 25, "listen Port.")
	flag.StringVar(&listen, "listen", ":25", "listen IP Addr and Port.")
	flag.Parse()

	if len(listen) < 1 {
		listen = fmt.Sprintf("%s:%d", host, port)
	}

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	conf := relaymail.Config{
		Listen:  listen,
		NextMTA: flag.Args()[0],
	}

	relaymail.ListenAndServe(conf)

}

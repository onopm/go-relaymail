package main

import (
	relaymail "github.com/onopm/go-relaymail"
)

func main() {

	relaymail.ListenAndServe(":10025")

}

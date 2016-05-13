package relaymail

import (
	"fmt"
	"log"

	"github.com/labstack/gommon/color"
)

func infof(format string, a ...interface{}) {
	log.Print(color.Green(fmt.Sprintf(format, a...)))
}

func warnf(format string, a ...interface{}) {
	log.Print(color.Yellow(fmt.Sprintf(format, a...)))
}

func critf(format string, a ...interface{}) {
	log.Print(color.Red(fmt.Sprintf(format, a...)))
}

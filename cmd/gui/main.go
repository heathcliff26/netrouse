package main

import (
	"flag"

	"github.com/heathcliff26/netrouse/pkg/gui"
)

func main() {
	// TODO: Silence fyne Error messages
	flag.Parse()
	initializeLogger()

	gui.New().Run()
}

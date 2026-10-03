package main

import (
	"flag"

	"github.com/heathcliff26/netrouse/pkg/gui"
)

func main() {
	flag.Parse()
	initializeLogger()

	gui.New().Run()
}

package main

import (
	"os"

	"github.com/LimEJET/mumbler"
)

func main() {
	m := mumbler.New()

	m.Name("Mumbler")
	m.Server("localhost")
	m.AddTracks(os.Args[1:]...)
	m.Connect()
	m.Play()

}

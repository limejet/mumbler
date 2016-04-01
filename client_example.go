package main

import (
	"os"

	"github.com/LimEJET/mumbler"
)

func main() {
	m := mumbler.New()

	m.Name("Mumbler")
	m.Server("localhost")
	m.AddFile(os.Args[1:]...)
	m.Connect()
	m.MoveToChannel("MP3 Player", true) // Creates it if nonexistent
	m.Play()
	m.Disconnect()

}

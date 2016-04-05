package main

import (
	"flag"
	"fmt"

	"github.com/LimEJET/mumbler"
)

var (
	NAME        = flag.String("n", "mumbler", "name to join server with")
	SERVER      = flag.String("s", "localhost", "server to join")
	CHANNEL     = flag.String("c", "Root", "channel to join")
	AVCONV      = flag.Bool("avconv", false, "Use avconv instead of ffmpeg")
	CERT        = flag.String("cert", "", "TLS certificate file (PEM)")
	SKIP_VERIFY = flag.Bool("skip-verify", false, "Skip TLS certificate verification")
	LOOP        = flag.Bool("loop", false, "Loop playlist")
	VOLUME      = flag.Float64("v", 1.0, "set max volume")
)

func main() {
	flag.Parse()
	m := mumbler.New()
	m.Name(*NAME)
	m.Server(*SERVER)
	m.Repeat(*LOOP)
	m.Volume(float32(*VOLUME))
	m.AudioDucking(0.90)

	if *AVCONV {
		m.Command("avconv")
	}

	if *CERT != "" {
		// no key file support for the example.
		m.Certificate(*CERT, "")
		m.SetTLSInsecureSkipVerify(*SKIP_VERIFY)
	} else {
		// we don't care that the certificate we didn't supply is invalid
		m.SetTLSInsecureSkipVerify(true)
	}

	m.AddFile(flag.Args()...)

	if err := m.Connect(); err != nil {
		fmt.Println(err)
		return
	}

	m.MoveToChannel(*CHANNEL, true)

	if err := m.Play(); err != nil {
		fmt.Println(err)
		return
	}

	if err := m.Disconnect(); err != nil {
		fmt.Println(err)
	}
}

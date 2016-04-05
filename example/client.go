package main

import (
	"flag"
	"fmt"
	"os"

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
)

func main() {
	flag.Parse()
	m := mumbler.New()
	m.Name(*NAME)
	m.Server(*SERVER)
	m.Repeat(false)
	m.AudioDucking(0.80)

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

	m.AddFile(os.Args[1:]...)

	if err := m.Connect(); err != nil {
		fmt.Println(err)
		return
	}

	m.MoveToChannel(*CHANNEL, true)
	if err := m.Play(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(m.Disconnect())
}

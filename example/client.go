package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/LimEJET/mumbler"
)

var (
	NAME    = flag.String("n", "mumbler", "name to join server with")
	SERVER  = flag.String("s", "localhost", "server to join")
	CHANNEL = flag.String("c", "Root", "channel to join")
	CREATE  = flag.Bool("create", false, "Allow create channel")
	AVCONV  = flag.Bool("avconv", false, "Use avconv instead of ffmpeg")
)

func main() {
	flag.Parse()
	m := mumbler.New()
	m.Name(*NAME)
	m.Server(*SERVER)
	if *AVCONV {
		m.Command("avconv")
	}
	
	m.AddFile(os.Args[1:]...)
	
	if err: m.Connect(); err != nil {
		fmt.Println(err)
		return
  }
	
	m.MoveToChannel(*CHANNEL, *CREATE)
	if err := m.Play(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(m.Disconnect())
}

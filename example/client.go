package main

import (
        "flag"
        "fmt"
        "github.com/LimEJET/mumbler"
        "os"
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
        err := m.Connect()
        if err != nil {
                fmt.Println(err)
                return
        }
        m.MoveToChannel(*CHANNEL, *CREATE)
        err = m.Play()
        if err != nil {
                fmt.Println(err)
                return
        }

        fmt.Println(m.Disconnect())
}

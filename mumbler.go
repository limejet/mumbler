package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleffmpeg"
	_ "github.com/layeh/gumble/opus"
)

var (
	SERVER   = flag.String("server", "localhost:64738", "Murmur server address")
	USERNAME = flag.String("user", "mumbler", "client username")
	INSECURE = flag.Bool("insecure", false, "skip server certificate verification")
	CERT     = flag.String("cert", "", "user certificate file (PEM)")
	KEY      = flag.String("key", "", "certificate key")
	CHANNEL  = flag.String("chan", "Root", "Channel to join on connect")
	CREATE   = flag.Bool("create-chan", false, "Create channel if nonexistent. Only works if -chan is specified")
	FILE     = flag.String("f", "", "mp3 file to play")
)

func main() {
	// stolen from gumbleutils
	flag.Parse()

	host, port, err := net.SplitHostPort(*SERVER)
	if err != nil {
		host = *SERVER
		port = strconv.Itoa(gumble.DefaultPort)
	}

	// client
	config := gumble.NewConfig()
	config.Username = *USERNAME
	config.Address = net.JoinHostPort(host, port)

	loadCert(config)

	client := gumble.NewClient(config)

	//keepAlive := make(chan bool)

	/*client.Attach(gumbleutil.Listener{
		Disconnect: func(e *gumble.DisconnectEvent) {
			keepAlive <- true
		},
	})
	*/
	if err := client.Connect(); err != nil {
		fmt.Printf("%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}
	defer client.Disconnect()

	// wait until connected
	for client.State() == gumble.StateConnecting {
	}
	moveToChannel(client)

	source := gumbleffmpeg.SourceFile(*FILE)
	stream := gumbleffmpeg.New(client, source)

	if err := stream.Play(); err != nil {
		fmt.Println("Play error:", err)
		os.Exit(1)
	}
	stream.Wait()

	// <-keepAlive
}

func loadCert(config *gumble.Config) {
	if *INSECURE {
		config.TLSConfig.InsecureSkipVerify = true
	}
	if *CERT != "" {
		if *KEY == "" {
			KEY = CERT
		}
		cert, err := tls.LoadX509KeyPair(*CERT, *KEY)
		if err != nil {
			fmt.Printf("Failed to read certificate: %v\n", err)
			os.Exit(1)
		}

		config.TLSConfig.Certificates = append(config.TLSConfig.Certificates, cert)
	}
}

func moveToChannel(client *gumble.Client) {
	if *CHANNEL != "Root" {
		connChan := client.Channels.Find(*CHANNEL)
		if connChan == nil {
			root := client.Channels.Find()
			if (*root.Permission() & gumble.PermissionMakeTemporaryChannel) == 0 {
				fmt.Println("Fatal: cannot create channels in Root")
				os.Exit(1)
			}

			if !*CREATE {
				fmt.Printf("Fatal: Channel \"%v\" does not exist!\n", *CHANNEL)
				os.Exit(1)
			}
			root.Add(*CHANNEL, true)
			connChan := client.Channels.Find(*CHANNEL)
			if connChan == nil {
				fmt.Printf("Fatal: Failed to create channel %v\n", *CHANNEL)
				os.Exit(1)
			}
			fmt.Printf("Created \"%v\"\n", *CHANNEL)
		}
		client.Self.Move(connChan)
		fmt.Printf("Moved to \"%v\"\n", *CHANNEL)
	}
}

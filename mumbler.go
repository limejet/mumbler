package mumbler

import (
	"crypto/tls"
	"errors"
	"net"
	"strconv"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleffmpeg"
	_ "github.com/layeh/gumble/opus"
)

type Mumbler struct {
	playlist []string
	config   *gumble.Config
	client   *gumble.Client
}

func New() *Mumbler {
	config := gumble.NewConfig()
	config.TLSConfig.InsecureSkipVerify = true
	return &Mumbler{
		config: config,
	}
}

func (m *Mumbler) Name(n string) {
	m.config.Username = n
}

func (m *Mumbler) Password(n string) {
	m.config.Password = n
}

func (m *Mumbler) MoveToChannel(name string, create bool) error {
	target := m.client.Channels.Find(name)

	if target == nil {
		root := m.client.Channels.Find()
		if !create {
			return errors.New("Nonexistent channel " + name)
		}
		if (*root.Permission() & gumble.PermissionMakeTemporaryChannel) == 0 {
			return errors.New("Permission error: Cannot create channels in root")
		}

		root.Add(name, true)

		target := m.client.Channels.Find(name)
		if target == nil {
			return errors.New("Failed to create channel " + name)
		}
	}
	m.client.Self.Move(target)
	return nil

}

func (m *Mumbler) AddFile(file ...string) {
	m.playlist = append(m.playlist, file...)
}

func (m *Mumbler) ClearFiles() {
	m.playlist = []string{}
}

func (m *Mumbler) Certificate(file, keyfile string) error {
	m.config.TLSConfig.InsecureSkipVerify = false

	if keyfile == "" {
		keyfile = file
	}
	cert, err := tls.LoadX509KeyPair(file, keyfile)
	if err != nil {
		return err
	}

	m.config.TLSConfig.Certificates = append(m.config.TLSConfig.Certificates, cert)
	return nil
}

func (m *Mumbler) Play() error {
	for _, file := range m.playlist {
		source := gumbleffmpeg.SourceFile(file)
		stream := gumbleffmpeg.New(client, source)

		if err := stream.Play(); err != nil {
			return err
		}
	}
}

func (m *Mumbler) Server(address string) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		host = address
		port = strconv.Itoa(gumble.DefaultPort)
	}

	m.config.Address = net.JoinHostPort(host, port)
}

func (m *Mumbler) Connect() error {

	m.client = gumble.NewClient(m.config)

	if err := m.client.Connect(); err != nil {
		return err
	}

	// wait until connected
	for m.client.State() == gumble.StateConnecting {
	}

	return nil
}

func (m *Mumbler) Disconnect() error {
	return m.client.Disconnect()
}

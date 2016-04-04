package mumbler

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleffmpeg"
	_ "github.com/layeh/gumble/opus"
)

type Mumbler struct {
	playlist     []Source
	config       *gumble.Config
	client       *gumble.Client
	stream       *gumbleffmpeg.Stream
	command      string
	loop         bool
	audioDucking bool
	volume       float32
	duckVolume   float32
}

func New() *Mumbler {
	config := gumble.NewConfig()
	return &Mumbler{config: config}
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
	for _, v := range file {
		m.playlist = append(m.playlist, NewFileSource(v))
	}
}

func (m *Mumbler) AddReader(reader ...io.Reader) {
	for _, v := range reader {
		m.playlist = append(m.playlist, NewReaderSource(v))
	}
}

func (m *Mumbler) AddReadCloser(reader ...io.ReadCloser) {
	for _, v := range reader {
		m.playlist = append(m.playlist, NewReadCloserSource(v))
	}
}

func (m *Mumbler) ClearPlaylist() {
	m.playlist = []Source{}
}

func (m *Mumbler) Command(c string) {
	m.command = c
}

func (m *Mumbler) Certificate(file, keyfile string) error {
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

func (m *Mumbler) SetTLSInsecureSkipVerify(b bool) {
	m.config.TLSConfig.InsecureSkipVerify = b
}

func (m *Mumbler) Play() error {
	for {
		for _, playlistItem := range m.playlist {
			source := playlistItem.GetSource()
			m.stream = gumbleffmpeg.New(m.client, source)
			m.volume = m.stream.Volume

			if m.command != "" {
				m.stream.Command = m.command
			}
			if err := m.stream.Play(); err != nil {
				return err
			}
			m.stream.Wait()
		}
		if !m.loop {
			break
		}
	}
	return nil
}

func (m *Mumbler) Repeat(l bool) {
	m.loop = l
}

func (m *Mumbler) AudioDucking(volume float32) {
	if volume != 0.0 {
		m.audioDucking = true
		m.duckVolume = 1 - volume
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
	if m.audioDucking {
		m.client.AttachAudio(m)
	}
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

func (m *Mumbler) OnAudioStream(e *gumble.AudioStreamEvent) {
	//BEWARE: here be hax
	go func() {
		t := time.Tick(500 * time.Millisecond)
		currVol := m.volume
		for {
			select {
			case <-t:
			volDn:
				for ; currVol < m.volume; currVol += 0.1 {
					time.Sleep(50 * time.Millisecond) // for a nice fade
					m.stream.Volume = currVol
					select {
					case <-e.C:
						break volDn
					default:
					}
				}
			case <-e.C:
				m.stream.Volume = m.duckVolume
				currVol = m.duckVolume
			}
		}
	}()
}

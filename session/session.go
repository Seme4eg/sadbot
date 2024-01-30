package session

import (
	"github.com/bwmarrin/discordgo"
	"github.com/seme4eg/sadbot/stream"
)

type Session struct {
	Session *discordgo.Session
	prefix  string
	// map of all streams that contains one stream per one guild
	Streams *stream.Streams
}

func new(token, prefix string) (*Session, error) {
	ses, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}
	return &Session{
		Session: ses,
		prefix:  prefix,
		Streams: &stream.Streams{List: make(map[string]*stream.Stream)},
	}, nil
}

func (s *Session) addHandlers() {
	s.Session.AddHandler(s.ready)         // ready events.
	s.Session.AddHandler(s.messageCreate) // messageCreate events.
	s.Session.AddHandler(s.guildCreate)   // guildCreate events.
}

// create websocket connection with discord
func (s *Session) open() error {
	return s.Session.Open()
}

func (s *Session) Close() error {
	return s.Session.Close()
}

func OpenSession(token, prefix string) (*Session, error) {
	session, err := new(token, prefix)
	if err != nil {
		return nil, err
	}

	session.addHandlers()

	if err := session.open(); err != nil {
		return nil, err
	}

	return session, nil
}

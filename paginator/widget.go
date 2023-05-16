package paginator

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// WidgetHandler ...
type WidgetHandler func(*Widget, *discordgo.MessageReaction)

// Widget is a message embed with reactions for buttons.
// Accepts custom handlers for reactions.
type Widget struct {
	Embed     *discordgo.MessageEmbed
	Message   *discordgo.Message
	Ses       *discordgo.Session
	ChannelID string
	Close     chan bool

	Handlers map[string]WidgetHandler
	// stores the handlers keys in the order they were added
	Keys []string
}

// NewWidget returns a pointer to a Widget object
func NewWidget(ses *discordgo.Session, channelID string) *Widget {
	return &Widget{
		ChannelID: channelID,
		Ses:       ses,
		Keys:      []string{},
		Handlers:  map[string]WidgetHandler{},
		Close:     make(chan bool),
	}
}

// Spawn spawns the widget in channel w.ChannelID
func (w *Widget) Spawn() error {
	// Create initial message.
	msg, err := w.Ses.ChannelMessageSendEmbed(w.ChannelID, w.Embed)
	if err != nil {
		return err
	}
	w.Message = msg

	// Add reaction buttons
	for _, v := range w.Keys {
		w.Ses.MessageReactionAdd(w.Message.ChannelID, w.Message.ID, v)
	}

	w.cleanupOnMessageDelete()

	var reaction *discordgo.MessageReaction
	for {
		select {
		case k := <-nextMessageReactionAddC(w.Ses):
			reaction = k.MessageReaction
		case <-w.Close:
			return nil
		}

		// Ignore reactions sent by bot
		if reaction.MessageID != w.Message.ID || w.Ses.State.User.ID == reaction.UserID {
			continue
		}

		if v, ok := w.Handlers[reaction.Emoji.Name]; ok {
			go v(w, reaction)
		}

		// delete reactions after they were added
		go func() {
			time.Sleep(time.Millisecond * 250)
			w.Ses.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID)
		}()
	}
}

// adds a handler for the given emoji name
func (w *Widget) Handle(emojiName string, handler WidgetHandler) {
	if _, ok := w.Handlers[emojiName]; !ok {
		w.Keys = append(w.Keys, emojiName)
		w.Handlers[emojiName] = handler
	}
}

// updated original message with new embed
func (w *Widget) UpdateEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return w.Ses.ChannelMessageEditEmbed(w.ChannelID, w.Message.ID, embed)
}

// sends close signal to Close channel of the widget on message delete
func (w *Widget) cleanupOnMessageDelete() {
	w.Ses.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageDelete) {
		// if current widget message gets deleted - send true to widgets' close chan
		if e.ID == w.Message.ID {
			w.Close <- true
		}
	})
}

// NextMessageReactionAddC returns a channel for the next MessageReactionAdd event
func nextMessageReactionAddC(s *discordgo.Session) chan *discordgo.MessageReactionAdd {
	out := make(chan *discordgo.MessageReactionAdd)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageReactionAdd) {
		out <- e
	})
	return out
}

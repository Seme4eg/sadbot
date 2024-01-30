package paginator

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Paginator provides a method for creating a navigatable embed
type Paginator struct {
	mu        sync.Mutex // dont expose mutex methods
	Ses       *discordgo.Session
	ChannelID string
	Close     chan bool

	Embed   *discordgo.MessageEmbed
	Message *discordgo.Message
	Pages   []*discordgo.MessageEmbed
	Index   int

	// Loop back to the beginning or end when on the first or last page.
	Loop bool

	Handlers map[string]func()
	// stores the handlers keys in the order they were added
	Keys []string
}

// NewPaginator returns a new Paginator
func NewPaginator(ses *discordgo.Session, channelID string) *Paginator {
	p := &Paginator{
		Pages:     []*discordgo.MessageEmbed{},
		ChannelID: channelID,
		Ses:       ses,
		Keys:      []string{},
		Handlers:  map[string]func(){},
		Close:     make(chan bool),
	}
	p.addReactions()

	return p
}

// Handle adds a handler for the given emoji name
func (p *Paginator) Handle(emojiName string, handler func()) {
	if _, ok := p.Handlers[emojiName]; !ok {
		p.Keys = append(p.Keys, emojiName)
		p.Handlers[emojiName] = handler
	}
}

func (p *Paginator) addReactions() {
	p.Handle(NavStart, func() { p.Goto(0) })
	p.Handle(NavLeft, func() { p.Goto(p.Index - 1) })
	p.Handle(NavRight, func() { p.Goto(p.Index + 1) })
	p.Handle(NavEnd, func() { p.Goto(len(p.Pages) - 1) })
}

// Spawn spawns the paginator in channel p.ChannelID
func (p *Paginator) Spawn(index ...int) error {
	// Sets the footers of all added pages to their page numbers.
	p.SetPageFooters()

	// set initial index of a paginator
	if index[0] != 0 {
		p.Index = index[0]
	}
	p.Embed = p.Pages[index[0]]

	// Create initial message.
	msg, err := p.Ses.ChannelMessageSendEmbed(p.ChannelID, p.Embed)
	if err != nil {
		return err
	}
	p.Message = msg

	// Add reaction buttons
	for _, v := range p.Keys {
		p.Ses.MessageReactionAdd(p.Message.ChannelID, p.Message.ID, v)
	}

	// remove all event handlers on current paginator message delete event
	p.cleanupOnMessageDelete()

	var reaction *discordgo.MessageReaction
	for {
		select {
		case k := <-nextMessageReactionAddC(p.Ses):
			reaction = k.MessageReaction
		case <-p.Close:
			return nil
		}

		// Ignore reactions sent by bot
		if reaction.MessageID != p.Message.ID || p.Ses.State.User.ID == reaction.UserID {
			continue
		}

		// call handler for given emoji name if found
		if v, ok := p.Handlers[reaction.Emoji.Name]; ok {
			go v()
		}

		// delete reactions after they were added
		go func() {
			time.Sleep(time.Millisecond * 250)
			p.Ses.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID)
		}()
	}
}

// Add adds a page to the paginator
func (p *Paginator) Add(embeds ...*discordgo.MessageEmbed) {
	p.Pages = append(p.Pages, embeds...)
}

// Goto jumps to the requested page index
func (p *Paginator) Goto(index int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	l := len(p.Pages)

	if index < 0 || index >= l {
		// Set the queue back to the beginning if Loop is enabled.
		if p.Loop {
			if index < 0 {
				index = l - 1 + index
			}
			if index >= l {
				index -= l
			}
		} else {
			fmt.Println("Tried to update paginator with index out of bounds")
			return
		}
	}

	p.Index = index

	// updated original message with new embed
	_, err := p.Ses.ChannelMessageEditEmbed(p.ChannelID, p.Message.ID, p.Pages[p.Index])
	if err != nil {
		fmt.Println("Error updating embed:", err)
	}
}

// SetPageFooters sets the footer of each embed to
// Be its page number out of the total length of the embeds.
// need to call it in spawn function cuz only then we know total page amount
func (p *Paginator) SetPageFooters() {
	for index, embed := range p.Pages {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("#[%d / %d]", index+1, len(p.Pages)),
		}
	}
}

// cleanupOnMessageDelete sends close signal to Close channel of the paginator
// on message delete
func (p *Paginator) cleanupOnMessageDelete() {
	p.Ses.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageDelete) {
		// if current paginator message gets deleted - send true to it's close chan
		if e.ID == p.Message.ID {
			p.Close <- true
		}
	})
}

// nextMessageReactionAddC returns a channel for the next MessageReactionAdd event
func nextMessageReactionAddC(s *discordgo.Session) chan *discordgo.MessageReactionAdd {
	out := make(chan *discordgo.MessageReactionAdd)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageReactionAdd) {
		out <- e
	})
	return out
}

package paginator

import (
	"errors"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var ErrIndexOutOfBounds = errors.New("err: Index is out of bounds")

// Paginator provides a method for creating a navigatable embed
type Paginator struct {
	sync.Mutex
	Pages []*discordgo.MessageEmbed
	Index int

	// Loop back to the beginning or end when on the first or last page.
	Loop   bool
	Widget *Widget

	lockToUser bool
}

// NewPaginator returns a new Paginator
func NewPaginator(ses *discordgo.Session, channelID string) *Paginator {
	p := &Paginator{
		Pages:  []*discordgo.MessageEmbed{},
		Widget: NewWidget(ses, channelID),
	}
	p.addReactions()

	return p
}

func (p *Paginator) addReactions() {
	p.Widget.Handle(NavBeginning, func(w *Widget, r *discordgo.MessageReaction) {
		p.Goto(0)
	})
	p.Widget.Handle(NavLeft, func(w *Widget, r *discordgo.MessageReaction) {
		if err := p.Prev(); err == nil {
			p.Widget.UpdateEmbed(p.Pages[p.Index])
		}
	})
	p.Widget.Handle(NavRight, func(w *Widget, r *discordgo.MessageReaction) {
		if err := p.Next(); err == nil {
			p.Widget.UpdateEmbed(p.Pages[p.Index])
		}
	})
	p.Widget.Handle(NavEnd, func(w *Widget, r *discordgo.MessageReaction) {
		p.Goto(len(p.Pages) - 1)
	})
}

// Spawn spawns the paginator in channel p.ChannelID
func (p *Paginator) Spawn() error {
	p.Widget.Embed = p.Pages[p.Index]

	return p.Widget.Spawn()
}

// Add a page to the paginator
func (p *Paginator) Add(embeds ...*discordgo.MessageEmbed) {
	p.Pages = append(p.Pages, embeds...)
}

// Next sets the page index to the next page
func (p *Paginator) Next() error {
	p.Lock()
	defer p.Unlock()

	if p.Index < 0 || p.Index+1 >= len(p.Pages) {
		// Set the queue back to the beginning if Loop is enabled.
		if p.Loop {
			p.Index = 0
			return nil
		}
		return ErrIndexOutOfBounds
	}

	p.Index++
	return nil
}

// Prev sets the current page index to the previous page.
func (p *Paginator) Prev() error {
	p.Lock()
	defer p.Unlock()

	if p.Index < 0 || p.Index+1 >= len(p.Pages) {
		// Set the queue back to the beginning if Loop is enabled.
		if p.Loop {
			p.Index = len(p.Pages) - 1
			return nil
		}
		return ErrIndexOutOfBounds
	}

	p.Index--
	return nil
}

// Goto jumps to the requested page index
func (p *Paginator) Goto(index int) {
	p.Lock()
	defer p.Unlock()
	p.Index = index
	_, err := p.Widget.UpdateEmbed(p.Pages[p.Index])
	if err != nil {
		fmt.Println("Error updating embed:", err)
	}
}

// SetPageFooters sets the footer of each embed to
// Be its page number out of the total length of the embeds.
func (p *Paginator) SetPageFooters() {
	for index, embed := range p.Pages {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("#[%d / %d]", index+1, len(p.Pages)),
		}
	}
}

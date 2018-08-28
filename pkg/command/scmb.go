package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chneau/anecdote/pkg/anecdote"
)

// SCMB ...
type SCMB struct{}

// Name ...
func (c *SCMB) Name() string {
	return "scmb"
}

// Description ...
func (c *SCMB) Description() string {
	return "> se coucher moins bete"
}

// MessageCreate ...
func (c *SCMB) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) error {
	anecdotes, err := anecdote.SCMB()
	if err != nil {
		return err
	}
	s.ChannelMessageSend(m.ChannelID, anecdotes[0].String())
	return nil
}

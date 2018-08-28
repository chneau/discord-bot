package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chneau/anecdote/pkg/anecdote"
)

// SI ...
type SI struct{}

// Name ...
func (c *SI) Name() string {
	return "si"
}

// Description ...
func (c *SI) Description() string {
	return "> savoir inutile"
}

// MessageCreate ...
func (c *SI) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) error {
	anecdotes, err := anecdote.SI()
	if err != nil {
		return err
	}
	s.ChannelMessageSend(m.ChannelID, anecdotes[0].String())
	return nil
}

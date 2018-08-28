package command

import (
	"github.com/bwmarrin/discordgo"
)

// Ping ...
type Ping struct{}

// Name ...
func (c *Ping) Name() string {
	return "ping"
}

// Description ...
func (c *Ping) Description() string {
	return "> pong"
}

// MessageCreate ...
func (c *Ping) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) error {
	s.ChannelMessageSend(m.ChannelID, "pong")
	return nil
}

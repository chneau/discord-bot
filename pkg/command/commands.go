package command

import (
	"github.com/bwmarrin/discordgo"
)

// Commander ...
type Commander interface {
	Name() string
	Description() string
	MessageCreate(*discordgo.Session, *discordgo.MessageCreate) error
}

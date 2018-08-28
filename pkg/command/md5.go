package command

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// MD5 ...
type MD5 struct{}

// Name ...
func (c *MD5) Name() string {
	return "md5"
}

// Description ...
func (c *MD5) Description() string {
	return "__text__ > to hash text"
}

// MessageCreate ...
func (c *MD5) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) error {
	strTok := strings.SplitN(m.Content, " ", 2)
	if len(strTok) > 1 {
		hasher := md5.New()
		hasher.Write([]byte(strTok[1])) // Recast the string to a byte array
		result := hex.EncodeToString(hasher.Sum(nil))
		s.ChannelMessageSend(m.ChannelID, result)
	}
	return nil
}

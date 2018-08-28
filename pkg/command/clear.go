package command

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// Clear ...
type Clear struct{}

// Name ...
func (c *Clear) Name() string {
	return "clear"
}

// Description ...
func (c *Clear) Description() string {
	return "> clear all messages"
}

// MessageCreate ...
func (c *Clear) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) error {
	// TODO check that user is either Owner or have the rights to delete messages !
	messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
	if err != nil {
		return err
	}
	messagesID := []string{}
	fourteenDaysAgo := time.Now()
	fourteenDaysAgo = fourteenDaysAgo.Add(-time.Hour * 24 * 14)
	for i := range messages {
		m := messages[i]
		t, _ := m.Timestamp.Parse()
		if t.After(fourteenDaysAgo) {
			messagesID = append(messagesID, m.ID)
		} else {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}
	}
	if err := s.ChannelMessagesBulkDelete(m.ChannelID, messagesID); err != nil {
		return err
	}
	return nil
}

package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/chneau/anecdote/pkg/anecdote"

	"github.com/bwmarrin/discordgo"
)

func init() {
	log.SetPrefix("[BOT] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// checkError
func ce(err error, msg string) {
	if err != nil {
		log.Panicln(msg, err)
	}
}

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		panic("No DISCORD_BOT_TOKEN env.")
	}
	discord, err := discordgo.New("Bot " + token)
	ce(err, "discordgo.New")
	defer discord.Close()
	discord.AddHandler(messageCreate)
	err = discord.Open()
	ce(err, "discord.Open")
	println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	println()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.Println(m.Author.Username, m.Content)
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, ".") {
		defer s.ChannelMessageDelete(m.ChannelID, m.ID)
	} else {
		return
	}
	strTok := strings.SplitN(m.Content, " ", 2)
	switch strTok[0] {
	case ".help":
		s.ChannelMessageSend(m.ChannelID, help)
	case ".md5":
		if len(strTok) > 1 {
			s.ChannelMessageSend(m.ChannelID, MD5Hash(strTok[1]))
		}
	case ".ping":
		date, err := m.Timestamp.Parse()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
			break
		}
		now := time.Now()
		s.ChannelMessageSend(m.ChannelID, "pong "+now.Sub(date).String())
	case ".annecdote":
		fallthrough
	case ".scmb":
		anecdotes, err := anecdote.SCMB()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
			break
		}
		s.ChannelMessageSend(m.ChannelID, anecdotes[0].String())
	case ".si":
		anecdotes, err := anecdote.SI()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
			break
		}
		s.ChannelMessageSend(m.ChannelID, anecdotes[0].String())
	case ".rmrf":
		// TODO check that user is either Owner or have the rights to delete messages !
		messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
			break
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
			s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		}
	}
}

// MD5Hash ...
func MD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text)) // Recast the string to a byte array
	return hex.EncodeToString(hasher.Sum(nil))
}

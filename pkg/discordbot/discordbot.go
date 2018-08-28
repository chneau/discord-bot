package discordbot

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/chneau/discord-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
)

// DiscordBot ...
type DiscordBot struct {
	Session  *discordgo.Session
	Commands map[string]command.Commander
	Prefix   string
}

// Close ...
func (d *DiscordBot) Close() error {
	return d.Session.Close()
}

func (d *DiscordBot) help() string {
	help := "List of commands:\n"
	for i := range d.Commands {
		help += d.Prefix + d.Commands[i].Name() + " " + d.Commands[i].Description() + "\n"
	}
	return help
}

func (d *DiscordBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, d.Prefix) {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	} else {
		return
	}
	cmd := strings.SplitN(m.Content, " ", 2)[0]
	cmd = cmd[1:]
	if val, ok := d.Commands[cmd]; ok {
		if err := val.MessageCreate(s, m); err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		}
		return
	}
	if cmd == "help" {
		s.ChannelMessageSend(m.ChannelID, d.help())
		return
	}
}

func (d *DiscordBot) signalHandler() {
	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
		d.Close()
	}()
}

// Add ...
func (d *DiscordBot) Add(cmd command.Commander) {
	d.Commands[cmd.Name()] = cmd
}

func (d *DiscordBot) init() {
	d.signalHandler()
	d.Session.AddHandler(d.messageCreate)
	if err := d.Session.Open(); err != nil {
		panic(err)
	}
}

// NewDefault ...
func NewDefault(token string) *DiscordBot {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	d := &DiscordBot{
		Session:  discord,
		Commands: make(map[string]command.Commander),
		Prefix:   ".",
	}
	d.init()
	d.Add(&command.MD5{})
	d.Add(&command.Ping{})
	d.Add(&command.SCMB{})
	d.Add(&command.SI{})
	d.Add(&command.Clear{})
	d.Add(&command.YT{})
	println("Bot is now running.  Press CTRL-C to exit.")
	return d
}

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chneau/discord-bot/pkg/discordbot"
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
	d := discordbot.NewDefault(os.Getenv("DISCORD_BOT_TOKEN"))
	defer d.Close()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	println()
}

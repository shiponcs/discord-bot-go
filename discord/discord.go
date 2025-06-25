package discord

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	// discordBotMentionRegexFmt supports also nicknames (the exclamation mark).
	// Read more: https://discordjs.guide/miscellaneous/parsing-mention-arguments.html#how-discord-mentions-work
	DiscordBotMentionRegexFmt = "^<@!?%s>"

	// discordMaxMessageSize max size before a message should be uploaded as a file.
	discordMaxMessageSize = 2000
)

type Discord struct {
	Token   string
	Session *discordgo.Session
}

func NewDiscord(token string, intents discordgo.Intent, stateEnabled bool) (*Discord, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	dg.Identify.Intents = intents
	dg.StateEnabled = stateEnabled

	discord := &Discord{
		Token:   token,
		Session: dg,
	}

	return discord, nil
}

type MessageHandlerType interface{}

func (d *Discord) AddHandler(handler MessageHandlerType) {
	d.Session.AddHandler(handler)
}

func (d *Discord) Init() {
	err := d.Session.Open()
	defer d.Session.Close()
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}
	fmt.Println("Bot is now running. Press Ctrl+C to exit.")

	// Wait for exit
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("Bot shutting down.")
}

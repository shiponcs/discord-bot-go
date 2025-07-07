package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/shiponcs/bot-go/discord"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable is required")
	}

	err := startBot(token)
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	fmt.Println("Exiting the App")
}

func startBot(token string) error {
	discordObj, err := discord.NewDiscord(token, discordgo.IntentsAll, true)
	if err != nil {
		log.Printf("Error creating Discord session: %v", err)
		return err
	}

	discordObj.AddHandler(discord.MessageHandler)
	discordObj.Init() // blocking call

	return nil
}

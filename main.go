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

	discordObj, err := discord.NewDiscord(token, discordgo.IntentsAll, true)
	if err != nil {
		fmt.Printf("ERROR NewDiscord %v\n", err)
		return
	}

	discordObj.AddHandler(discord.MessageHandler)
	discordObj.Init() // blocking call

	fmt.Println("Exiting the App")
}

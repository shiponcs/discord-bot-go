package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	webhookURL = os.Getenv("N8N_WEBHOOK_URL")
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable is required")
	}

	if webhookURL == "" {
		webhookURL = "http://103.203.237.218:5876/webhook-test/26c7308d-7072-4fe3-a8a3-5abbc5982013"

	}

	dg, err := discordgo.New("Bot " + token)
	dg.Identify.Intents = discordgo.IntentsAll
	dg.StateEnabled = true

	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	dg.AddHandler(messageHandler)

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}
	defer dg.Close()

	fmt.Println("Bot is now running. Press Ctrl+C to exit.")

	// Wait for exit
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("Bot shutting down.")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Only respond to direct messages
	// TODO(shiponcs): state cache doesn't work for DMs
	// it fails for DMs and then we call the channel directly which
	// makes a API request to fetch the information. We need
	// find a way around this.
	channel, err := s.State.Channel(m.ChannelID)
	if err == discordgo.ErrStateNotFound {
		channel, err = s.Channel(m.ChannelID)
		if err != nil {
			log.Printf("Error fetching channel: %v", err)
			return
		}
	}

	if channel.Type != discordgo.ChannelTypeDM {
		fmt.Println("Not a DM, ignore")
		return
	}

	payload := map[string]interface{}{
		"username":  m.Author.Username + "#" + m.Author.Discriminator,
		"user_id":   m.Author.ID,
		"content":   m.Content,
		"channel":   "Direct Message",
		"timestamp": m.Timestamp,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error encoding payload: %v", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending to n8n: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Message sent to n8n: %d", resp.StatusCode)
}

package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var (
	webhookURL = os.Getenv("N8N_WEBHOOK_URL")
)

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	content := m.Content
	channel_id := ""
	channel_name := "Direct message"

	if channel.Type != discordgo.ChannelTypeDM {
		fmt.Println("Not a DM")
		botMentionRegex, err := discordBotMentionRegex(s.State.User.ID)
		if err != nil {
			fmt.Println("can't build botMentionRegex")
			return
		}

		req, found := findAndTrimBotMention(botMentionRegex, m.Content)
		if !found {
			fmt.Println("Not a DM nor a mentioned one, ignoring")
			return
		}
		content = req
		channel_id = m.ChannelID
		channel_name = "channel message"
		fmt.Println("the mentioned message ", req)
	}

	payload := map[string]interface{}{
		"username":  m.Author.Username + "#" + m.Author.Discriminator,
		"user_id":   m.Author.ID,
		"content":   content,
		"channel":   channel_name,
		"timestamp": m.Timestamp,
	}
	if channel_id != "" {
		payload["channel_id"] = channel_id
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error encoding payload: %v", err)
		return
	}

	if webhookURL == "" {
		webhookURL = "http://103.203.237.218:5876/webhook-test/26c7308d-7072-4fe3-a8a3-5abbc5982013"

	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending to n8n: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Message sent to n8n: %d", resp.StatusCode)
}

func discordBotMentionRegex(botID string) (*regexp.Regexp, error) {
	botMentionRegex, err := regexp.Compile(fmt.Sprintf(DiscordBotMentionRegexFmt, botID))
	if err != nil {
		return nil, fmt.Errorf("while compiling bot mention regex: %w", err)
	}

	return botMentionRegex, nil
}

func findAndTrimBotMention(botMentionRegex *regexp.Regexp, msg string) (string, bool) {
	if !botMentionRegex.MatchString(msg) {
		return "", false
	}

	return botMentionRegex.ReplaceAllString(msg, ""), true
}

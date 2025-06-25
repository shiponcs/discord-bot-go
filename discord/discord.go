package discord

const (
	// discordBotMentionRegexFmt supports also nicknames (the exclamation mark).
	// Read more: https://discordjs.guide/miscellaneous/parsing-mention-arguments.html#how-discord-mentions-work
	DiscordBotMentionRegexFmt = "^<@!?%s>"

	// discordMaxMessageSize max size before a message should be uploaded as a file.
	discordMaxMessageSize = 2000
)

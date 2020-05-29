package handler

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// check if the message is "!airhorn"
	if strings.HasPrefix(m.Content, ".") {
		if m.Content == ".ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}
	}
}

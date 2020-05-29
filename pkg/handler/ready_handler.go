package handler

import (
	"github.com/bwmarrin/discordgo"
)

func ReadyHandler(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, ".help")
}

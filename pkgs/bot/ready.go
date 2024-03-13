package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ready(s *discordgo.Session, _ *discordgo.Ready) {
	log.Println("connected to discord as", s.State.User.Username)
}

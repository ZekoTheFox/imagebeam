package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func ready(s *discordgo.Session, e *discordgo.Ready) {
	log.Println("successfully connected to discord as", s.State.User.Username)
}

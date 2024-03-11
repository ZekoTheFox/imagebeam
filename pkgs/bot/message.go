package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func messageCreate(_ *discordgo.Session, e *discordgo.MessageCreate) {
	if len(e.Attachments) < 1 {
		return
	}

	// only the first attachment will ever be used
	attachment := e.Attachments[0]

	log.Println("attachment found /", attachment.ContentType, ",", attachment.Filename, ",", attachment.URL)
}

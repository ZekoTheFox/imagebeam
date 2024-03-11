package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"zeko.party/imagebeam/pkgs/webapi"
)

func checkList(ids []string, target string) bool {
	// special case: if ids list is empty, default to true
	if len(ids) == 0 {
		return true
	}

	for _, id := range ids {
		if id == target {
			return true
		}
	}

	return false
}

func messageCreate(_ *discordgo.Session, e *discordgo.MessageCreate) {
	if len(e.Attachments) < 1 || e.Author.Bot {
		return
	}

	if !checkList(config.PermittedChannels, e.ChannelID) {
		return
	}

	if !checkList(config.PermittedUsers, e.Author.ID) {
		return
	}

	// only the first attachment will ever be used/shown
	attachment := e.Attachments[0]

	validContentType := []string{
		"image/png",
		"image/jpeg",
		"image/webp",
		"image/gif",
		"image/avif",
	}
	validType := false
	for _, fileType := range validContentType {
		if attachment.ContentType == fileType {
			validType = true
		}
	}

	if !validType {
		return
	}

	log.Println("attachment found from user", e.Author.Username, "("+e.Author.ID+")")
	log.Println("- content-type: " + attachment.ContentType)
	log.Println("- filename: " + attachment.Filename)

	webapi.Images <- webapi.Image{
		Url: attachment.URL,
	}
}

package bot

import (
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

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

var (
	regexDiscordCDN   = regexp.MustCompile(`https:\/\/cdn\.discordapp\.com\/attachments\/[0-9]{18,19}\/[0-9]{18,19}\/.+\.(png|jpg|jpeg|webp|gif|avif)\?ex=[a-z0-9]{8}&is=[a-z0-9]{8}&hm=[a-z0-9]{64}&?`)
	regexDiscordMedia = regexp.MustCompile(`https:\/\/media\.discordapp\.net\/attachments\/[0-9]{18,19}\/[0-9]{18,19}\/.+\.(png|jpg|jpeg|webp|gif|avif)\?ex=[a-z0-9]{8}&is=[a-z0-9]{8}&hm=[a-z0-9]{64}(&=&format=webp(&quality=lossless)?&width=[0-9]{1,4}&height=[0-9]{1,4})?&?`)
	regexTenor        = regexp.MustCompile(`https:\/\/tenor\.com\/view\/[\w\-(%\d{2})]+-[0-9]+`)
	regexTenorMedia   = regexp.MustCompile(`https:\/\/media1.tenor\.com\/m\/[\w\d]+\/[\w\-]+\.gif`)
)

func queueImage(url string) {
	webapi.Images <- webapi.Image{
		Url: url,
	}
}

func handleAttachment(e *discordgo.MessageCreate) {

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

	queueImage(attachment.ProxyURL)
}

func handleLinks(e *discordgo.MessageCreate) {
	// check if there are valid links/embeds if no attachments were found
	message := e.Content
	matchedCDN := regexDiscordCDN.FindString(message)
	matchedMediaProxy := regexDiscordMedia.FindString(message)
	if matchedCDN != "" {
		// mutate cdn url into media proxy to optimize loading
		converted := strings.ReplaceAll(matchedCDN, "cdn.discordapp.com", "media.discordapp.net")
		if !strings.Contains(converted, ".gif") {
			converted = converted + "=&format=webp"
		}

		log.Println("link found from user", e.Author.Username, "("+e.Author.ID+")")

		queueImage(converted)
		return
	}

	if matchedMediaProxy != "" {
		log.Println("link found from user", e.Author.Username, "("+e.Author.ID+")")

		queueImage(matchedMediaProxy)
		return
	}

	// tenor is weird and doesn't have an easy way to query direct gif links
	matchedTenor := regexTenor.FindString(message)
	if matchedTenor != "" {
		log.Println("link found from user", e.Author.Username, "("+e.Author.ID+")")

		page, err := http.Get(matchedTenor)
		if err != nil {
			log.Println("warning: failed to crawl tenor webpage")
			return
		}

		pageText, err := io.ReadAll(page.Body)
		if err != nil {
			log.Println("warning: failed to read out tenor page data")
			return
		}

		text := string(pageText[:])

		resolvedMedia := regexTenorMedia.FindString(text)

		log.Println("resolved tenor link:", resolvedMedia)

		queueImage(resolvedMedia)
	}
}

func messageCreate(_ *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.Bot {
		return
	}

	if !checkList(config.PermittedChannels, e.ChannelID) {
		return
	}

	if !checkList(config.PermittedUsers, e.Author.ID) {
		return
	}

	// see if there are attachments
	if len(e.Attachments) >= 1 {
		handleAttachment(e)
		return
	}

	if regexDiscordCDN.MatchString(e.Content) || regexDiscordMedia.MatchString(e.Content) || regexTenor.MatchString(e.Content) {
		handleLinks(e)
		return
	}
}

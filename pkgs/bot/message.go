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
	regexDiscordCDN   = regexp.MustCompile(`https:\/\/cdn\.discordapp\.com\/attachments\/[0-9]{18,19}\/[0-9]{18,19}\/.+\.(png|jpg|jpeg|webp|gif)\?ex=[a-z0-9]{8}&is=[a-z0-9]{8}&hm=[a-z0-9]{64}&?`)
	regexDiscordMedia = regexp.MustCompile(`https:\/\/media\.discordapp\.net\/attachments\/[0-9]{18,19}\/[0-9]{18,19}\/[a-zA-Z0-9_\-\.]+\.(png|jpg|jpeg|webp|gif)\?ex=[0-9a-f]{8}&is=[0-9a-f]{8}&hm=[0-9a-f]{64}(&=&format=webp(&quality=lossless)?&width=[0-9]{1,4}&height=[0-9]{1,4})?&?`)
	regexDiscordEmoji = regexp.MustCompile(`https:\/\/cdn\.discordapp\.com\/emojis\/[0-9]{18,19}\.webp`)
	regexTenor        = regexp.MustCompile(`https:\/\/tenor\.com\/view\/([\w\-]|%[0-9]{2})+-[0-9]+`)
	regexTenorMedia   = regexp.MustCompile(`https:\/\/media1\.tenor\.com\/m\/[\w\d]+\/[\w\-]+\.gif`)
)

func queueImage(url string) {
	log.Println("processing link =>", strings.TrimPrefix(url, "https://")[:36]+"...")

	webapi.Images <- webapi.Image{
		Url: url,
	}
}

func handleAttachment(e *discordgo.MessageCreate) {
	// only the first attachment will ever be used/shown
	attachment := e.Attachments[0]

	filetype := []string{
		"image/png",
		"image/jpeg",
		"image/webp",
		"image/gif",
	}
	valid := false
	for _, entry := range filetype {
		if attachment.ContentType == entry {
			valid = true
		}
	}

	if !valid {
		return
	}

	log.Println("attachment found from user", e.Author.Username, "("+e.Author.ID+")")
	log.Println("- content-type: " + attachment.ContentType)
	log.Println("- filename: " + attachment.Filename)

	queueImage(attachment.ProxyURL)
}

func handleLinks(e *discordgo.MessageCreate) {
	message := e.Content
	matched := regexDiscordCDN.FindString(message)
	if matched != "" {
		// change cdn into media proxy url to optimize loading
		converted := strings.ReplaceAll(matched, "cdn.discordapp.com", "media.discordapp.net")
		if !strings.Contains(converted, ".gif") {
			converted = converted + "=&format=webp"
		}

		queueImage(converted)
		return
	}

	matched = regexDiscordEmoji.FindString(message)
	if matched != "" {
		queueImage(matched)
		return
	}

	matched = regexDiscordMedia.FindString(message)
	if matched != "" {
		queueImage(matched)
		return
	}

	// tenor is weird and buries their direct gif links
	matched = regexTenor.FindString(message)
	if matched != "" {
		page, err := http.Get(matched)
		if err != nil {
			log.Println("warning: failed to crawl tenor webpage")
			return
		}

		pageText, err := io.ReadAll(page.Body)
		if err != nil {
			log.Println("warning: failed to read out tenor page data")
			return
		}

		resolvedMedia := regexTenorMedia.FindString(string(pageText[:]))
		log.Println("resolved tenor link:", resolvedMedia)

		queueImage(resolvedMedia)
	}
}

func messageCreate(_ *discordgo.Session, e *discordgo.MessageCreate) {
	if e.Author.Bot || !checkList(config.PermittedChannels, e.ChannelID) || !checkList(config.PermittedUsers, e.Author.ID) {
		return
	}

	if len(e.Attachments) >= 1 {
		handleAttachment(e)
		return
	}

	text := e.Content
	if strings.HasPrefix(text, "https://") {
		handleLinks(e)
		return
	}

	if len(e.StickerItems) >= 0 && text == "" {
		sticker := e.StickerItems[0]
		extension := ""

		switch sticker.FormatType {
		case discordgo.StickerFormatTypeGIF:
			extension = ".gif"
		case discordgo.StickerFormatTypePNG,
			discordgo.StickerFormatTypeAPNG:
			extension = ".png"
		default:
			return
		}

		if len(extension) == 0 {
			return
		}

		queueImage("https://media.discordapp.net/stickers/" + sticker.ID + extension + "?size=320")
		return
	}

	emojis := e.GetCustomEmojis()
	if len(emojis) > 0 {
		emoji := emojis[0]

		queueImage("https://cdn.discordapp.com/emojis/" + emoji.ID + ".webp?size=128&quality=lossless")
	}
}

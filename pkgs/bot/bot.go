package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type BotConfig struct {
	Token             string
	PermittedChannels []string
	PermittedUsers    []string
}

func StartDiscordBot(config BotConfig) (session *discordgo.Session, err error) {
	bot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatal("error initializing discord bot", err)
		return nil, err
	}

	bot.AddHandler(ready)
	bot.AddHandler(messageCreate)

	bot.Identify.Intents = discordgo.IntentGuildMessages

	err = bot.Open()
	if err != nil {
		log.Fatal("error connecting to discord")
		return nil, err
	}

	log.Println("discord bot successfully started")
	return bot, nil
}
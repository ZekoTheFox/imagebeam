package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"zeko.party/imagebeam/pkgs/bot"
	"zeko.party/imagebeam/pkgs/webapi"
)

var (
	Token             string
	PermittedChannels string
	PermittedUsers    string
	Port              int
)

func init() {
	flag.StringVar(&Token, "token", "", "Discord bot authentication token (required)")
	flag.StringVar(&PermittedChannels, "channels", "", "Channel IDs which will be listened to and processed; assumes the bot has been given permission to the specified channels.\nAccepts a single channel ID, or multiple IDs separated by a single comma.\n(defaults to all available channels)")
	flag.StringVar(&PermittedUsers, "users", "", "User IDs which will be the sole users processed. Accepts a single user ID, or multiple IDs separated by a single comma. (defaults to all users)")
	flag.IntVar(&Port, "port", 8440, "The port the web server will use. May break if changed from default.")
	flag.Parse()
}

func validateSnowflake(id string) bool {
	matched, err := regexp.MatchString(id, "[0-9]{18,19}")
	return matched && err == nil
}

func processIDList(input string) []string {
	ids := []string{}

	if strings.Contains(input, ",") {
		for _, id := range strings.Split(input, ",") {
			if !validateSnowflake(id) {
				log.Println("warning: found invalid id ->", id)
				continue
			}

			ids = append(ids, id)
		}

		return ids
	}

	if !validateSnowflake(input) {
		log.Println("warning: found invalid id ->", input)
		return []string{} // default to empty array
	}

	ids = append(ids, input)
	return ids
}

func main() {
	log.Println("imagebeam is starting up...")

	log.Println("token = " + strings.Split(Token, ".")[0] + "...") // the first segment is just the bot's id

	if len(PermittedChannels) <= 17 {
		log.Println("warning: no permitted channels set; all available channels will be valid channels")
	}
	if len(PermittedUsers) <= 17 {
		log.Println("warning: no permitted users set; all images in permitted channels will be shown")
	}

	channels := processIDList(PermittedChannels)
	users := processIDList(PermittedUsers)

	botConfig := bot.BotConfig{
		Token:             Token,
		PermittedChannels: channels,
		PermittedUsers:    users,
	}

	bot, err := bot.StartDiscordBot(botConfig)
	if err != nil {
		log.Fatal("error attempting to start discord bot; check logs for more detail")
		return
	}

	defer bot.Close()

	webapi.StartWebAPI(Port)

	log.Println("imagebeam started")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig

	log.Println("shutting down...")
}

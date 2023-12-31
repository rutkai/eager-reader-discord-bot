package discord

import (
	"EagerReaderDiscordBot/ai"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"os"
	"regexp"
	"strings"
)

var blacklist = []string{
	"youtube.com",
	"youtu.be",
	"9gag.com",
}

func StartBot() {
	discord, err := discordgo.New("Bot " + getBotToken())
	if err != nil {
		log.Fatal().Err(err).Msg("error creating Discord session")
	}

	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.IntentsGuildMessages

	err = discord.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("error opening connection")
	}
	defer discord.Close()
}

func getBotToken() string {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		panic("BOT_TOKEN environment variable is missing")
	}
	return token
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Author.Bot {
		log.Debug().Str("Author username", m.Author.Username).Msg("Filtering out bot message")
		return
	}

	url := parseUrl(m.Content)
	if url == "" {
		log.Debug().Str("Message", m.Content).Msg("Message does not contain URL")
		return
	}
	log.Debug().Str("Url", url).Msg("URL found in message")

	if isBlacklisted(url) {
		log.Debug().Str("Url", url).Msg("URL is blacklisted")
		return
	}

	summary, err := ai.GetSummary(url, "")
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, "Something went wrong with generating the summary.")
		if err != nil {
			log.Error().Err(err).Msg("Message send failed")
		}
		return
	}

	log.Debug().Str("Response", summary).Msg("Sending response from OpenAI")
	_, err = s.ChannelMessageSend(m.ChannelID, summary)
}

func parseUrl(url string) string {
	re := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
	return re.FindString(url)
}

func isBlacklisted(url string) bool {
	for _, blacklistedUrl := range blacklist {
		if strings.Contains(url, blacklistedUrl) {
			return true
		}
	}
	return false
}

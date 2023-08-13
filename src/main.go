package main

import (
	"EagerReaderDiscordBot/ai"
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

func main() {
	initZerolog()

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

	log.Debug().Msg("Debug mode enabled!")
	log.Info().Msg("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Info().Msg("Bot received SIGTERM, terminating main process.")
	discord.Close()
}

func initZerolog() {
	debug := flag.Bool("debug", false, "sets log level to debug")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
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

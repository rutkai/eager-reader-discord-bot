package main

import (
	"EagerReaderDiscordBot/discord"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initZerolog()

	discord.SetAllowlist(getAllowlist())
	discord.SetBlocklist(getBlocklist())
	discord.StartBot()

	log.Debug().Msg("Debug mode enabled!")
	log.Info().Msg("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Info().Msg("Bot received SIGTERM, terminating main process.")
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

func getBlocklist() []string {
	return getConfigList("blocklist")
}

func getAllowlist() []string {
	return getConfigList("allowlist")
}

func getConfigList(config string) []string {
	fileContent, err := os.ReadFile(fmt.Sprintf("config/%s.json", config))
	if err != nil {
		return []string{}
	}

	var list []string
	err = json.Unmarshal(fileContent, &list)
	if err != nil {
		return []string{}
	}

	return list
}

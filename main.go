package main

import (
	"os"
	"os/signal"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/kang-makes/discord-chatbot/bot"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	// Flag and environment parsing.
	pflag.StringP("log-level", "v", "info", "Log level (verbosity)")
	pflag.StringP("config", "c", "config.yaml", "YAML were to put configurations fo servers and functions")

	pflag.Parse()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.SetEnvPrefix("DC")
	viper.AutomaticEnv()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err.Error())
	}
	logLevel := viper.GetString("log-level")
	configFile := viper.GetString("config")

	// Logging
	lvl, _ := log.ParseLevel(logLevel)
	log.SetLevel(lvl)

	log.Infof("loaded options: log-level:'%q', config:'%q'", logLevel, configFile)

	// Bot Start
	if fInfo, err := os.Stat(configFile); os.IsNotExist(err) || fInfo.IsDir() {
		log.Fatalf("config provided is not a file")
	}

	osSignal := make(chan os.Signal, 10)
	signal.Notify(osSignal, os.Interrupt)
	defer close(osSignal)

	log.Trace("creating bot")
	bot, err := bot.ParseConfig(configFile)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Info("launching bot")
	err = bot.Run(osSignal)
	if err != nil {
		log.Errorf("bot cannot continue: %v", err)
	}
}

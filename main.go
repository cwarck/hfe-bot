package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	telebot "gopkg.in/telebot.v4"

	"hfe-go/pkg/config"
	"hfe-go/pkg/statemachine"
	"hfe-go/pkg/telegram"
)

func main() {
	cfg := config.MustNewAppConfig()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.IsDebug() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	bot, err := telebot.NewBot(
		telebot.Settings{
			Token: cfg.Telegram.Token,
			Poller: &telebot.Webhook{
				AllowedUpdates: []string{"message", "callback_query"},
				Endpoint: &telebot.WebhookEndpoint{
					PublicURL: cfg.WebhookUrl,
				},
				Listen: ":" + cfg.Port,
			},
			Verbose: cfg.IsDebug(),
		},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to initialize the bot")
	}

	stateManager := statemachine.NewManager()

	telegram.SetMiddlewares(cfg, bot)
	telegram.SetHandlers(bot, cfg, stateManager)

	log.Info().Msg("start listening for updates")
	bot.Start()
}

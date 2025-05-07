package telegram

import (
	"hfe-go/pkg/config"

	telebot "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

func SetMiddlewares(cfg *config.AppConfig, bot *telebot.Bot) {
	if cfg.IsDebug() {
		bot.Use(middleware.Logger())
	}
	if len(cfg.Telegram.AllowedChatIds) > 0 {
		bot.Use(middleware.Whitelist(cfg.Telegram.AllowedChatIds...))
	}
}

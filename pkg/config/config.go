package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	Categories             []Category
	GoogleSheets           GoogleSheets
	Telegram               Telegram
	OpenExchangeRatesAppId string
	DefaultCurrency        string
	LogLevel               string
	Port                   string
	WebhookUrl             string
}

type Category struct {
	Name  string
	Emoji string
}

type Telegram struct {
	AllowedChatIds []int64
	Token          string
}

type GoogleSheets struct {
	SpreadsheetId string
	SheetName     string
}

func (c *AppConfig) IsDebug() bool {
	return c.LogLevel == "debug"
}

// MustGetEnv gets the environment variable. Panics if the environment variable is not set.
func mustGetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic(fmt.Sprintf("%s is a required environment variable", key))
}

// GetEnvOrDefault gets the environment variable or the default value if the environment variable is not set.
func getEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// MustNewAppConfig creates a new AppConfig. Panics on error.
func MustNewAppConfig() *AppConfig {
	// TODO: get categories from google sheets
	categories := []Category{
		{Name: "Bills", Emoji: "💸"},
		{Name: "Childcare", Emoji: "👶"},
		{Name: "Clothes", Emoji: "👚"},
		{Name: "Eating out", Emoji: "🍴"},
		{Name: "Education", Emoji: "🎓"},
		{Name: "Food delivery", Emoji: "🍔"},
		{Name: "Groceries", Emoji: "🥦"},
		{Name: "Healthcare", Emoji: "🏥"},
		{Name: "Hobbies", Emoji: "🎨"},
		{Name: "Other", Emoji: "❔"},
		{Name: "Rent", Emoji: "🏠"},
		{Name: "Self-care", Emoji: "💅"},
		{Name: "Shopping", Emoji: "🛍️"},
		{Name: "Subscriptions", Emoji: "✨"},
		{Name: "Transport", Emoji: "🚕"},
		{Name: "Travel", Emoji: "✈️"},
	}

	var allowedChatIds []int64
	for chatId := range strings.SplitSeq(mustGetEnv("TELEGRAM_ALLOWED_CHAT_IDS"), ",") {
		chatIdInt, err := strconv.ParseInt(chatId, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("invalid chatId in TELEGRAM_ALLOWED_CHAT_IDS: %s", chatId))
		}
		allowedChatIds = append(allowedChatIds, chatIdInt)
	}

	return &AppConfig{
		Categories: categories,
		Telegram: Telegram{
			Token:          mustGetEnv("TELEGRAM_TOKEN"),
			AllowedChatIds: allowedChatIds,
		},
		GoogleSheets: GoogleSheets{
			SpreadsheetId: mustGetEnv("GOOGLE_SHEETS_SPREADSHEET_ID"),
			SheetName:     mustGetEnv("GOOGLE_SHEETS_SHEET_NAME"),
		},
		OpenExchangeRatesAppId: mustGetEnv("OPENEXCHANGERATES_APP_ID"),
		DefaultCurrency:        mustGetEnv("DEFAULT_CURRENCY"),
		LogLevel:               getEnvOrDefault("LOG_LEVEL", "info"),
		Port:                   getEnvOrDefault("PORT", "3000"),
		WebhookUrl:             mustGetEnv("WEBHOOK_URL"),
	}
}

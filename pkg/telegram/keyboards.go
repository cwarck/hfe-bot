package telegram

import (
	"fmt"
	"hfe-go/pkg/config"

	telebot "gopkg.in/telebot.v4"
)

func MakeCategoriesButtons(categories []config.Category) []telebot.Btn {
	var buttons []telebot.Btn
	keyboard := &telebot.ReplyMarkup{}

	for i, category := range categories {
		unique := fmt.Sprintf("category_btn_%d", i)
		displayText := category.Name
		if category.Emoji != "" {
			displayText = fmt.Sprintf("%s %s", category.Emoji, category.Name)
		}
		buttons = append(buttons, keyboard.Data(displayText, unique, category.Name))
	}

	return buttons
}

func MakeCategoriesInlineKeyboard(categories []config.Category) *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}
	buttons := MakeCategoriesButtons(categories)
	columns := 3

	var rows []telebot.Row
	var offset int
	for i := 0; i < len(buttons); i += columns {
		offset += columns
		if offset > len(buttons) {
			offset = len(buttons)
		}
		rows = append(rows, keyboard.Row(buttons[i:offset]...))
	}
	keyboard.Inline(rows...)

	return keyboard
}

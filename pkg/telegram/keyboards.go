package telegram

import (
	"fmt"
	"strings"

	telebot "gopkg.in/telebot.v4"

	"hfe-go/pkg/config"
)

var (
	ButtonAddCategory = telebot.Btn{
		Unique: "add_category",
		Text:   "Add a category",
		Data:   "add_category",
	}
	ButtonRemoveCategory = telebot.Btn{
		Unique: "remove_category",
		Text:   "Remove a category",
		Data:   "remove_category",
	}
)

func CategoriesButtons(categories []config.Category) []telebot.Btn {
	var buttons []telebot.Btn
	for _, category := range categories {
		buttons = append(buttons, ButtonSelectCategory(category.Name, category.Emoji))
	}
	return buttons
}

func ButtonSelectCategory(name string, emoji string) telebot.Btn {
	unique := strings.ReplaceAll(name, "_", " ")
	return telebot.Btn{
		Unique: "category_" + unique,
		Text:   fmt.Sprintf("%s %s", emoji, name),
		Data:   name,
	}
}

func NewInlineKeyboard(buttons []telebot.Btn, columns int) *telebot.ReplyMarkup {
	if columns == 0 {
		columns = 1
	}

	keyboard := &telebot.ReplyMarkup{}
	var rows []telebot.Row
	if columns == 1 {
		for _, button := range buttons {
			rows = append(rows, keyboard.Row(button))
		}
	} else {
		var offset int
		for i := 0; i < len(buttons); i += columns {
			offset += columns
			if offset > len(buttons) {
				offset = len(buttons)
			}
			rows = append(rows, keyboard.Row(buttons[i:offset]...))
		}
	}
	keyboard.Inline(rows...)

	return keyboard
}

func NewCategoriesKeyboard(categories []config.Category) *telebot.ReplyMarkup {
	return NewInlineKeyboard(CategoriesButtons(categories), 3)
}

func NewAdminKeyboard() *telebot.ReplyMarkup {
	buttons := []telebot.Btn{
		ButtonAddCategory,
		ButtonRemoveCategory,
	}
	return NewInlineKeyboard(buttons, 1)
}

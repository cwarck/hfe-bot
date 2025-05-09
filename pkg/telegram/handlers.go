package telegram

import (
	"context"

	telebot "gopkg.in/telebot.v4"

	"hfe-go/pkg/config"
	"hfe-go/pkg/currency"
	"hfe-go/pkg/statemachine"
)

// HandleText handles the text message from the user.
func HandleText(cfg *config.AppConfig, stateManager *statemachine.Manager) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		ctx := context.Background()
		text := c.Text()
		senderId := c.Sender().ID

		// TODO: handle random text

		state, exists := stateManager.Get(senderId)
		if !exists {
			// Initialize expense state by default
			state = statemachine.NewAddExpenseState()
			stateManager.Set(senderId, state)
		}

		switch state.Current() {
		case statemachine.StateAddExpenseIdle:
			matches := currency.AmountRegexp.FindStringSubmatch(text)
			if matches != nil {
				if err := state.Event(ctx, statemachine.EventAddExpenseStart); err != nil {
					return err
				}
				if err := state.Event(ctx, statemachine.EventAddExpenseAmountEntered, text); err != nil {
					return err
				}
				return c.Send("Choose a category 👇", NewCategoriesKeyboard(cfg.Categories))
			}
			return nil

		case statemachine.StateAddExpenseWaitingComment:
			if err := state.Event(ctx, statemachine.EventAddExpenseCommentEntered, text, cfg); err != nil {
				return err
			}
			if err := state.Event(ctx, statemachine.EventAddExpenseSaveCompleted); err != nil {
				return err
			}
			stateManager.Delete(senderId)
			return c.Send("Saved successfully 👌")

		case statemachine.StateAdminAddCategoryWaitingCategory:
			if err := state.Event(ctx, statemachine.EventAdminAddCategoryCategoryEntered, text); err != nil {
				return err
			}
			stateManager.Delete(senderId)
			return c.Send("Saved successfully 👌")
		}

		return nil
	}
}

// HandleCategoriesKeyboardCallback handles the callback from the categories keyboard.
func HandleCategoriesKeyboardCallback(cfg *config.AppConfig, stateManager *statemachine.Manager) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		category := c.Callback().Data
		senderId := c.Sender().ID

		state, exists := stateManager.Get(senderId)
		if !exists || !state.Is(statemachine.StateAddExpenseWaitingCategory) {
			return c.RespondAlert("Enter amount first!")
		}

		err := state.Event(context.Background(), statemachine.EventAddExpenseCategorySelected, category)
		if err != nil {
			return err
		}

		return c.Edit("Ok, now you can send me a comment to add to the expense 💬 or a link to the receipt 🧾")
	}
}

// HandleCancel handles the cancel button.
func HandleCancel(stateManager *statemachine.Manager) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		senderId := c.Sender().ID
		state, exists := stateManager.Get(senderId)
		if !exists {
			return c.Send("No operation in progress")
		}

		if err := state.Event(context.Background(), statemachine.EventCancel); err != nil {
			return err
		}
		stateManager.Delete(senderId)

		return c.Send("Operation has been cancelled")
	}
}

func SetHandlers(bot *telebot.Bot, cfg *config.AppConfig, stateManager *statemachine.Manager) {
	// commands
	bot.Handle("/start", func(c telebot.Context) error { return c.Send("Cao! 👋") })
	bot.Handle("/cancel", HandleCancel(stateManager))

	// raw text messages
	bot.Handle(telebot.OnText, HandleText(cfg, stateManager))

	// categories keyboard
	for _, btn := range CategoriesButtons(cfg.Categories) {
		bot.Handle(&btn, HandleCategoriesKeyboardCallback(cfg, stateManager))
	}
}

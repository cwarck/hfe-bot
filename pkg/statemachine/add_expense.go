package statemachine

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"

	"hfe-go/pkg/config"
	"hfe-go/pkg/currency"
	"hfe-go/pkg/spreadsheets"
)

const (
	EventAddExpenseStart            = "event.add_expense.start"
	EventAddExpenseAmountEntered    = "event.add_expense.amount_entered"
	EventAddExpenseCategorySelected = "event.add_expense.category_selected"
	EventAddExpenseCommentEntered   = "event.add_expense.comment_entered"
	EventAddExpenseSaveCompleted    = "event.add_expense.save_completed"

	StateAddExpenseIdle            = "state.add_expense.idle"
	StateAddExpenseWaitingAmount   = "state.add_expense.waiting_amount"
	StateAddExpenseWaitingCategory = "state.add_expense.waiting_category"
	StateAddExpenseWaitingComment  = "state.add_expense.waiting_comment"
	StateAddExpenseReadyToSave     = "state.add_expense.ready_to_save"
)

func NewAddExpenseState() *fsm.FSM {
	events := fsm.Events{
		// Start expense entry
		{Name: EventAddExpenseStart, Src: []string{StateAddExpenseIdle}, Dst: StateAddExpenseWaitingAmount},

		// User has sent an amount, waiting for category
		{Name: EventAddExpenseAmountEntered, Src: []string{StateAddExpenseWaitingAmount}, Dst: StateAddExpenseWaitingCategory},

		// User has selected a category, waiting for comment
		{Name: EventAddExpenseCategorySelected, Src: []string{StateAddExpenseWaitingCategory}, Dst: StateAddExpenseWaitingComment},

		// User has sent a comment, ready to save
		{Name: EventAddExpenseCommentEntered, Src: []string{StateAddExpenseWaitingComment}, Dst: StateAddExpenseReadyToSave},

		// Save completed
		{Name: EventAddExpenseSaveCompleted, Src: []string{StateAddExpenseReadyToSave}, Dst: StateAddExpenseIdle},

		// Cancel at any point
		{Name: EventCancel, Src: []string{StateAddExpenseWaitingAmount, StateAddExpenseWaitingCategory, StateAddExpenseWaitingComment, StateAddExpenseReadyToSave}, Dst: StateAddExpenseIdle},
	}

	callbacks := fsm.Callbacks{
		"after_" + EventAddExpenseAmountEntered: func(_ context.Context, e *fsm.Event) {
			amount := e.Args[0].(string)
			e.FSM.SetMetadata("amount", amount)
		},

		"after_" + EventAddExpenseCategorySelected: func(_ context.Context, e *fsm.Event) {
			category := e.Args[0].(string)
			e.FSM.SetMetadata("category", category)
		},

		"after_" + EventAddExpenseCommentEntered: func(_ context.Context, e *fsm.Event) {
			comment := e.Args[0].(string)
			cfg := e.Args[1].(*config.AppConfig)

			amount, _ := e.FSM.Metadata("amount")
			category, _ := e.FSM.Metadata("category")

			// TODO: move to EventAddExpenseAmountEntered callback handling
			normalizedAmount, err := currency.GetAmount(amount.(string), cfg.OpenExchangeRatesAppId, cfg.DefaultCurrency)
			if err != nil {
				e.Cancel(fmt.Errorf("amount normalization error: %w", err))
				return
			}

			if err := spreadsheets.AddExpense(cfg.GoogleSheets, normalizedAmount, category.(string), comment); err != nil {
				e.Cancel(fmt.Errorf("spreadsheets.AddExpense error: %w", err))
				return
			}
		},
	}

	return fsm.NewFSM(StateAddExpenseIdle, events, callbacks)
}

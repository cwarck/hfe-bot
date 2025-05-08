package statemachine

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"
)

const (
	EventAdminAddCategoryStart           = "event.admin.add_category.start"
	EventAdminAddCategoryCategoryEntered = "event.admin.add_category.category_entered"

	StateAdminAddCategoryIdle            = "state.admin.add_category.idle"
	StateAdminAddCategoryWaitingCategory = "state.admin.add_category.waiting_category"
)

func NewAddCategoryState() *fsm.FSM {
	events := fsm.Events{
		// User is asked to send a category
		{Name: EventAdminAddCategoryStart, Src: []string{StateAdminAddCategoryIdle}, Dst: StateAdminAddCategoryWaitingCategory},

		// User has sent a category
		{Name: EventAdminAddCategoryCategoryEntered, Src: []string{StateAdminAddCategoryWaitingCategory}, Dst: StateAdminAddCategoryIdle},

		// Cancel at any point
		{Name: EventCancel, Src: []string{StateAdminAddCategoryWaitingCategory}, Dst: StateAdminAddCategoryIdle},
	}

	callbacks := fsm.Callbacks{
		"after_" + EventAdminAddCategoryCategoryEntered: func(_ context.Context, e *fsm.Event) {
			category := e.Args[0].(string)

			// TODO: save to google sheets

			fmt.Println("Category:", category)
		},
	}

	return fsm.NewFSM(StateAdminAddCategoryIdle, events, callbacks)
}

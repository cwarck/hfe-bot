package statemachine

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"
)

const (
	EventAdminRemoveCategoryStart           = "event.admin.remove_category.start"
	EventAdminRemoveCategoryCategoryEntered = "event.admin.remove_category.category_entered"

	StateAdminRemoveCategoryIdle            = "state.admin.remove_category.idle"
	StateAdminRemoveCategoryWaitingCategory = "state.admin.remove_category.waiting_category"
)

func NewRemoveCategoryState() *fsm.FSM {
	events := fsm.Events{
		// User is asked to send a category
		{Name: EventAdminRemoveCategoryStart, Src: []string{StateAdminRemoveCategoryIdle}, Dst: StateAdminRemoveCategoryWaitingCategory},

		// User has sent a category
		{Name: EventAdminRemoveCategoryCategoryEntered, Src: []string{StateAdminRemoveCategoryWaitingCategory}, Dst: StateAdminRemoveCategoryIdle},

		// Cancel at any point
		{Name: EventCancel, Src: []string{StateAdminRemoveCategoryWaitingCategory}, Dst: StateAdminRemoveCategoryIdle},
	}

	callbacks := fsm.Callbacks{
		"after_" + EventAdminRemoveCategoryCategoryEntered: func(_ context.Context, e *fsm.Event) {
			category := e.Args[0].(string)

			// TODO: remove from google sheets

			fmt.Println("Category:", category)
		},
	}

	return fsm.NewFSM(StateAdminRemoveCategoryIdle, events, callbacks)
}

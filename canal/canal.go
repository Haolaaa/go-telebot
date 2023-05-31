package canal

import (
	"telebot_v2/utils"

	"github.com/go-mysql-org/go-mysql/canal"
)

type EventHandler struct {
	canal.DummyEventHandler
}

func (h *EventHandler) OnRow(event *canal.RowsEvent) error {
	_, rows := GetParseValue(event)

	if _, ok := utils.InArray(event.Table.Name, []string{"video_release"}); ok {
		TableEventDispatcher(event, rows)
	}

	return nil
}

func (h *EventHandler) String() string {
	return "EventHandler"
}

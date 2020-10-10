package timebox

import (
	"github.com/kode4food/timebox/command"
	"github.com/kode4food/timebox/event"
	"github.com/kode4food/timebox/id"
	"github.com/kode4food/timebox/message"
	"github.com/kode4food/timebox/store"
	"github.com/kode4food/timebox/time"
)

// Top-Level Type Aliases
type (
	Command = command.Command
	Event   = event.Event
	ID      = id.ID
	Message = message.Message
	Store   = store.Store
	Time    = time.Time
)

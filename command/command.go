package command

import "github.com/kode4food/timebox/message"

type (
	// Command is a Message at the end of the day
	Command = message.Message

	// Handler for Commands
	Handler = message.Handler

	// Type of a Command
	Type = message.Type
)

// New will create a new Command instance
var New = message.New

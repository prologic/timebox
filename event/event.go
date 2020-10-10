package event

import "github.com/kode4food/timebox/message"

type (
	// Event is a Message at the end of the day
	Event = message.Message

	// Handler for Events
	Handler = message.Handler

	// List of Events
	List = message.List

	// Type of an Event
	Type = message.Type
)

var (
	// New will create a new Event
	New = message.New

	// EmptyList is an empty List of Events
	EmptyList = message.EmptyList
)

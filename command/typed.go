package command

import "fmt"

// Error messages
const (
	ErrTypeNotRegistered = "command type %s not registered with dispatcher"
)

// TypedHandler calls the Handler appropriate for a Command's Type
type TypedHandler map[Type]Handler

// Handler dispatches to an underlying registered Handler
func (t TypedHandler) Handler() Handler {
	return func(c *Command) error {
		if h, ok := t[c.Type]; ok {
			return h(c)
		}
		return fmt.Errorf(ErrTypeNotRegistered, c.Type)
	}
}

// Combine with other instances, yielding a new instance
func (t TypedHandler) Combine(handlers ...TypedHandler) TypedHandler {
	combined := append([]TypedHandler{t}, handlers...)
	return TypedHandlers(combined...)
}

// TypedHandlers unifies a TypedHandler set
func TypedHandlers(handlers ...TypedHandler) TypedHandler {
	res := TypedHandler{}
	for _, typed := range handlers {
		for commandType, handler := range typed {
			res[commandType] = handler
		}
	}
	return res
}

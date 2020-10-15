package server

import (
	"github.com/kode4food/timebox/event"
	"github.com/kode4food/timebox/store"
)

type resolver struct {
	store  store.Store
	source *event.Source
}

// NewResolver returns a new ResolverRoot instance that wraps our
// Store and an event sourcing interface for that Store
func NewResolver(s store.Store) ResolverRoot {
	return &resolver{
		store:  s,
		source: event.NewSource(s),
	}
}

// Mutation returns MutationResolver implementation.
func (r *resolver) Mutation() MutationResolver {
	return newMutationResolver(r)
}

// Query returns QueryResolver implementation.
func (r *resolver) Query() QueryResolver {
	return newQueryResolver(r)
}

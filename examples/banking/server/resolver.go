package server

import (
	"github.com/kode4food/timebox/event"
	"github.com/kode4food/timebox/store"
)

type resolver struct {
	store  store.Store
	source *event.Source
}

// NewResolver returns a new ResolverRoot instance
func NewResolver(s store.Store) ResolverRoot {
	r := &resolver{
		store:  s,
		source: event.NewSource(s),
	}
	return r
}

// Mutation returns MutationResolver implementation.
func (r *resolver) Mutation() MutationResolver {
	return newMutationResolver(r)
}

// Query returns QueryResolver implementation.
func (r *resolver) Query() QueryResolver {
	return newQueryResolver(r)
}

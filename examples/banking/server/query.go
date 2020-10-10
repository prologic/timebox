package server

import (
	"context"

	"banking/model"

	"github.com/kode4food/timebox/event"
	"github.com/kode4food/timebox/id"
	"github.com/kode4food/timebox/store"
)

type queries struct{ *resolver }

func newQueryResolver(r *resolver) QueryResolver {
	return &queries{
		resolver: r,
	}
}

func (q *queries) AccountStatus(
	_ context.Context, accountID id.ID) (*model.AccountStatus, error) {
	var res *model.AccountStatus
	err := q.source.With(accountID,
		func(a *event.Aggregate, result store.Result) error {
			acc, err := model.HydrateFrom(a, result)
			if err != nil {
				return err
			}
			res = &model.AccountStatus{
				AccountID: acc.AccountID,
				Balance:   acc.Balance,
			}
			return nil
		},
	)
	return res, err
}

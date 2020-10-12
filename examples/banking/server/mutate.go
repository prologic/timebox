package server

import (
	"context"

	"banking/model"

	"github.com/kode4food/timebox/command"
	"github.com/kode4food/timebox/id"
	"github.com/kode4food/timebox/message"
)

const retryCount = 10

type mutations struct {
	*resolver
	handle command.Handler
}

func newMutationResolver(r *resolver) MutationResolver {
	c := Handler(r.source)
	h := command.Retry(retryCount, c)
	return &mutations{
		resolver: r,
		handle:   h,
	}
}

func (m *mutations) OpenAccount(
	_ context.Context, input model.OpenAccount,
) (id.ID, error) {
	accountID := id.New()
	return accountID, m.invoke(OpenAccountCommand, OpenAccountWithID{
		OpenAccount: input,
		AccountID:   accountID,
	})
}

func (m *mutations) DepositMoney(
	_ context.Context, input model.TransferMoney,
) (id.ID, error) {
	return input.AccountID, m.invoke(DepositMoneyCommand, input)
}

func (m *mutations) WithdrawMoney(
	_ context.Context, input model.TransferMoney,
) (id.ID, error) {
	return input.AccountID, m.invoke(WithdrawMoneyCommand, input)
}

func (m *mutations) invoke(
	msgType message.Type, payload message.Payload,
) error {
	if err := m.handle(command.New(msgType, payload)); err != nil {
		return err
	}
	return nil
}

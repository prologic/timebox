package model

import (
	"github.com/kode4food/timebox"
	"github.com/kode4food/timebox/event"
	"github.com/kode4food/timebox/message"
	"github.com/kode4food/timebox/store"
)

// Event Names
const (
	AccountOpened  = "account-opened"
	MoneyDeposited = "account-amount-deposited"
	MoneyWithdrawn = "account-amount-withdrawn"
)

type Account struct {
	AccountID
	Owner   string
	Balance *Money
}

// TypedInstantiator is the instantiator for Events
var TypedInstantiator = message.TypedInstantiator{
	AccountOpened: func() message.Payload {
		return &AccountOpenedEvent{}
	},
	MoneyDeposited: func() message.Payload {
		return &MoneyDepositedEvent{}
	},
	MoneyWithdrawn: func() message.Payload {
		return &MoneyWithdrawnEvent{}
	},
}

// HydrateFrom creates a new Account instance and hydrates it with the
// Events that can be retrieved from the provided Store Result. Those
// Events are then applied the specified instance
func HydrateFrom(a *event.Aggregate, result store.Result) (*Account, error) {
	e, err := result.Events()
	if err != nil {
		return nil, err
	}
	acc := &Account{}
	a.ApplyTo(acc.Applier())
	a.HydrateFrom(e)
	return acc, nil
}

// Applier returns an Applier for the Account aggregate
func (a *Account) Applier() event.Applier {
	ta := event.TypedApplier{
		AccountOpened:  makeAccountOpened(a),
		MoneyDeposited: makeMoneyDeposited(a),
		MoneyWithdrawn: makeMoneyWithdrawn(a),
	}
	return ta.Applier()
}

func makeAccountOpened(a *Account) event.Applier {
	return func(e *timebox.Event) {
		p := e.Payload.(*AccountOpenedEvent)
		a.AccountID = p.AccountID
		a.Owner = p.Owner
		a.Balance = NewMoney(0, CurrencyEur)
	}
}

func makeMoneyDeposited(a *Account) event.Applier {
	return func(e *timebox.Event) {
		p := e.Payload.(*MoneyDepositedEvent)
		res, _ := a.Balance.Add(p.DepositedAmount)
		a.Balance = res
	}
}

func makeMoneyWithdrawn(a *Account) event.Applier {
	return func(e *timebox.Event) {
		p := e.Payload.(*MoneyWithdrawnEvent)
		res, _ := a.Balance.Subtract(p.WithdrawnAmount)
		a.Balance = res
	}
}

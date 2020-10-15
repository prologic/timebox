package model

import (
	"github.com/kode4food/timebox"
	"github.com/kode4food/timebox/event"
	"github.com/kode4food/timebox/message"
	"github.com/kode4food/timebox/store"
)

// Event Names
const (
	AccountOpenedEvent  = "account-opened"
	MoneyDepositedEvent = "account-amount-deposited"
	MoneyWithdrawnEvent = "account-amount-withdrawn"
)

type Account struct {
	AccountID
	Owner   string
	Balance *Money
}

// TypedInstantiator is the instantiator for Events
var TypedInstantiator = message.TypedInstantiator{
	AccountOpenedEvent: func() message.Payload {
		return &AccountOpened{}
	},
	MoneyDepositedEvent: func() message.Payload {
		return &MoneyDeposited{}
	},
	MoneyWithdrawnEvent: func() message.Payload {
		return &MoneyWithdrawn{}
	},
}

// HydrateFrom creates a new Account instance and hydrates it with the
// Events that can be retrieved from the provided Store Result. Those
// Events are then applied to the specified instance.
//
// This is the absolute most simple way to do this, but it means that
// instances are always built from the entire historical record. If
// aggregates have a lot of events, it may eventually make sense to
// perform snapshots.
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

// Applier returns an Applier for the Account aggregate. The Applier
// it returns is a TypedApplier, but could even be a single Handler
// function driven by a switch statement
func (a *Account) Applier() event.Applier {
	ta := event.TypedApplier{
		AccountOpenedEvent:  makeAccountOpened(a),
		MoneyDepositedEvent: makeMoneyDeposited(a),
		MoneyWithdrawnEvent: makeMoneyWithdrawn(a),
	}
	return ta.Applier()
}

func makeAccountOpened(a *Account) event.Applier {
	return func(e *timebox.Event) {
		p := e.Payload.(*AccountOpened)
		a.AccountID = p.AccountID
		a.Owner = p.Owner
		a.Balance = NewMoney(0, CurrencyEur)
	}
}

func makeMoneyDeposited(a *Account) event.Applier {
	return func(e *timebox.Event) {
		p := e.Payload.(*MoneyDeposited)
		res, _ := a.Balance.Add(p.DepositedAmount)
		a.Balance = res
	}
}

func makeMoneyWithdrawn(a *Account) event.Applier {
	return func(e *timebox.Event) {
		p := e.Payload.(*MoneyWithdrawn)
		res, _ := a.Balance.Subtract(p.WithdrawnAmount)
		a.Balance = res
	}
}

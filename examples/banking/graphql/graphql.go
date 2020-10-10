package graphql

//go:generate go run github.com/99designs/gqlgen

import (
	"encoding/json"
	"fmt"
	"io"

	"banking/model"

	"github.com/99designs/gqlgen/graphql"
	"github.com/kode4food/timebox"
	"github.com/kode4food/timebox/id"
)

const (
	errNotID     = "%T is not an ID"
	errNotObject = "%T is not an Object"
	errNotMoney  = "JSON value is not Money"
)

func MarshalTimeboxID(id timebox.ID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = fmt.Fprintf(w, `"%s"`, id.String())
	})
}

func UnmarshalTimeboxID(v interface{}) (timebox.ID, error) {
	switch v := v.(type) {
	case string:
		return id.Parse(v)
	default:
		return id.Nil, fmt.Errorf(errNotID, v)
	}
}

func MarshalMoney(m *model.Money) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		b, _ := json.Marshal(m)
		_, _ = w.Write(b)
	})
}

func UnmarshalMoney(v interface{}) (*model.Money, error) {
	switch v := v.(type) {
	case map[string]interface{}:
		var m model.Money
		return &m, UnmarshalMoneyObject(&m, v)
	default:
		return nil, fmt.Errorf(errNotObject, v)
	}
}

func UnmarshalMoneyObject(m *model.Money, obj map[string]interface{}) error {
	cents, okCents := obj["cents"]
	currency, okCurrency := obj["currency"]
	if !okCents || !okCurrency {
		return fmt.Errorf(errNotMoney)
	}
	if currency, ok := currency.(string); ok {
		if cents, ok := cents.(int64); ok {
			m.Currency = model.Currency(currency)
			m.Cents = model.Cents(cents)
		}
		return nil
	}
	return fmt.Errorf(errNotMoney)
}

package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoneyAdd(t *testing.T) {
	as := assert.New(t)

	m1 := NewMoney(10, EUR)
	m2 := NewMoney(34, EUR)
	m3 := NewMoney(61, USD)

	r1, err := m1.Add(m2)
	as.Nil(err)
	as.NotEqual(m1, r1)
	as.NotEqual(m2, r1)

	as.Equal(Cents(10), m1.Cents)
	as.Equal(Cents(44), r1.Cents)
	as.True(r1.IsPositive())
	as.Equal(EUR, r1.Currency)

	r2, err := m1.Add(m3)
	as.Nil(r2)
	as.Equal(err, fmt.Errorf(ErrIncompatibleCurrency, string(EUR), USD))
}

func TestMoneySubtract(t *testing.T) {
	as := assert.New(t)

	m1 := NewMoney(10, EUR)
	m2 := NewMoney(34, EUR)
	m3 := NewMoney(61, USD)

	r1, err := m1.Subtract(m2)
	as.Nil(err)
	as.NotEqual(m1, r1)
	as.NotEqual(m2, r1)

	as.Equal(Cents(10), m1.Cents)
	as.Equal(Cents(-24), r1.Cents)
	as.True(r1.IsNegative())
	as.Equal(EUR, r1.Currency)

	r2, err := m1.Subtract(m3)
	as.Nil(r2)
	as.Equal(err, fmt.Errorf(ErrIncompatibleCurrency, string(EUR), USD))
}

func TestMoneyEquality(t *testing.T) {
	as := assert.New(t)

	m1 := NewMoney(10, EUR)
	m2 := NewMoney(10, USD)
	m3 := NewMoney(10, EUR)

	as.False(m1.Equals(m2))
	as.True(m1.Equals(m3))

	r1, err := m1.Subtract(m3)
	as.Nil(err)
	as.NotEqual(m1, r1)
	as.NotEqual(m3, r1)
	as.True(r1.IsZero())
}

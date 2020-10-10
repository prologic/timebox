package model

import "fmt"

type (
	// Money stores monetary values associated with specific currencies
	Money struct {
		Cents    `json:"cents"`
		Currency `json:"currency"`
	}

	// Cents represents a monetary atom (ex: for USD, 100 Cents = 1 Dollar)
	Cents int
)

// Error messages
const (
	ErrIncompatibleCurrency   = "currencies (%s and %s) are incompatible"
	ErrPositiveFactorRequired = "positive factor required. got %d"
)

// Currency types
const (
	EUR Currency = "EUR"
	USD Currency = "USD"
)

// NewMoney creates Money out of thin air
func NewMoney(cents Cents, currency Currency) *Money {
	return &Money{Cents: cents, Currency: currency}
}

// Add one Money value to another, returning a new value
func (m *Money) Add(other *Money) (*Money, error) {
	if m.Currency != other.Currency {
		err := fmt.Errorf(ErrIncompatibleCurrency, m.Currency, other.Currency)
		return nil, err
	}
	return &Money{
		Cents:    m.Cents + other.Cents,
		Currency: m.Currency,
	}, nil
}

// Subtract one Money value from another, returning a new value
func (m *Money) Subtract(other *Money) (*Money, error) {
	if m.Currency != other.Currency {
		err := fmt.Errorf(ErrIncompatibleCurrency, m.Currency, other.Currency)
		return nil, err
	}
	return &Money{
		Cents:    m.Cents - other.Cents,
		Currency: m.Currency,
	}, nil
}

// IsPositive returns whether this value is a positive amount
func (m *Money) IsPositive() bool {
	return m.Cents > 0
}

// IsNegative returns whether this value is negative amount
func (m *Money) IsNegative() bool {
	return m.Cents < 0
}

// IsZero returns whether this value is zero
func (m *Money) IsZero() bool {
	return m.Cents == 0
}

// Equals checks two monetary values for equality
func (m *Money) Equals(other *Money) bool {
	return m.Cents == other.Cents && m.Currency == other.Currency
}

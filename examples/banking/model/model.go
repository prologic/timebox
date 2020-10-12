package model

import "fmt"

// Error messages
const (
	errPositiveAmount = "positive amount required"
)

// Checker is an interface that a Command might implement to validate
// its payload fields
type Checker interface {
	Check() error
}

func (c *TransferMoney) Check() error {
	if !c.Amount.IsPositive() {
		return fmt.Errorf(errPositiveAmount)
	}
	return nil
}

package types

import (
	"math/big"
	"time"
)

type Account struct {
	Name          string    // A slash separated hierarchy of words.
	Type          string    // BANK, EXPENSE, INCOME, ASSET, CASH.
	Unit          string    // A currency or security name.
	Denom         int       // The unit denominator (usually 100).
	Description   string    // A free text description.
	LastReconcile time.Time // The time of last reconciliation.
}

type Transaction struct {
	Date        time.Time // The value date of the transaction.
	Stamp       time.Time // When the transaction was entered.
	Description string
	Notes       string // Additional notes
	Number      string // A sequence number (checks...)
	Flows       []Flow
}

// A Flow is a part of a split transaction. A flow is positive for
// debit actions, negative for credit actions.
type Flow struct {
	Memo           string
	Account        *Account `json:"-"`
	Price          *big.Rat
	Reconciled     bool
	ReconciledTime time.Time
}

// An accounting book.
type Book struct {
	Accounts     []*Account
	Transactions []Transaction
}

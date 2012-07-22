package types

import (
	"strconv"
	"strings"
	"time"
)

// Cents represent a money amount in 1/100s.
type Cents int64

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

type Flow struct {
	Memo           string
	Account        *Account `json:"-"`
	Price          string
	Reconciled     bool
	ReconciledTime time.Time
}

// An accounting book.
type Book struct {
	Accounts     []*Account
	Transactions []Transaction
}

func ParsePrice(s string) float64 {
	if sl := strings.IndexRune(s, '/'); sl < 0 {
		x, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		return float64(x)
	} else {
		num, den := s[:sl], s[sl+1:]
		n, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		d, err := strconv.Atoi(den)
		if err != nil {
			panic(err)
		}
		return float64(n) / float64(d)
	}
	panic("unreachable")
}

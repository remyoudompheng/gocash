package types

import (
	"fmt"
	"math/big"
	"sort"
	"time"
)

type Account struct {
	Name          string     // A slash separated hierarchy of words.
	Type          string     // BANK, EXPENSE, INCOME, ASSET, CASH.
	Unit          string     // A currency or security name.
	Denom         int        // The unit denominator (usually 100).
	Description   string     // A free text description.
	LastReconcile time.Time  // The time of last reconciliation.
	Children      []*Account `json:"-"`
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
	Price          *Amount
	Reconciled     bool
	ReconciledTime time.Time
	Parent         *Transaction `json:"-"`
}

type Amount big.Rat

func (amt *Amount) Rat() *big.Rat               { return (*big.Rat)(amt) }
func (amt *Amount) SetRat(rat *big.Rat) *Amount { (*big.Rat)(amt).Set(rat); return amt }

func (amt *Amount) String() string {
	return (*big.Rat)(amt).FloatString(2)
}

func (amt *Amount) MarshalJSON() (s []byte, err error) {
	return []byte((*big.Rat)(amt).RatString()), nil
}

func (amt *Amount) UnmarshalJSON(s []byte) error {
	_, ok := (*big.Rat)(amt).SetString(string(s))
	if !ok {
		return fmt.Errorf("invalid price string %q", s)
	}
	return nil
}

func (x *Amount) Add(y *Amount) *Amount { (*big.Rat)(x).Add((*big.Rat)(x), (*big.Rat)(y)); return x }

// An accounting book.
type Book struct {
	Accounts     []*Account
	Transactions []Transaction

	// Computed data.
	Balance map[*Account]*Amount `json:"-"`
	// Flows by account.
	Flows map[*Account][]*Flow `json:"-"`
}

func (book *Book) Recompute() {
	book.Balance = book.computeBalances()
	book.Flows = book.sortFlows()
}

// BalanceCents returns the per-account balance 
func (book *Book) computeBalances() (balance map[*Account]*Amount) {
	balance = make(map[*Account]*Amount, len(book.Accounts))
	for _, act := range book.Accounts {
		balance[act] = new(Amount)
	}
	for _, trn := range book.Transactions {
		for _, f := range trn.Flows {
			bal := balance[f.Account]
			balance[f.Account] = bal.Add(f.Price)
		}
	}
	return
}

func (book *Book) sortFlows() (flows map[*Account][]*Flow) {
	flows = make(map[*Account][]*Flow, len(book.Accounts))
	for it, trn := range book.Transactions {
		for i, f := range trn.Flows {
			flows[f.Account] = append(flows[f.Account], &trn.Flows[i])
			trn.Flows[i].Parent = &book.Transactions[it]
		}
	}
	for _, actflows := range flows {
		sort.Sort(flowsByDate(actflows))
	}
	return
}

type flowsByDate []*Flow

func (s flowsByDate) Len() int           { return len(s) }
func (s flowsByDate) Less(i, j int) bool { return s[i].Parent.Date.Before(s[j].Parent.Date) }
func (s flowsByDate) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

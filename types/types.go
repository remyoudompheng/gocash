package types

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"time"
)

type Account struct {
	Id            GUID
	Name          string     // A slash separated hierarchy of words.
	Type          string     // BANK, EXPENSE, INCOME, ASSET, CASH.
	Unit          string     // A currency or security name.
	Denom         int        // The unit denominator (usually 100).
	Description   string     // A free text description.
	LastReconcile time.Time  // The time of last reconciliation.
	Children      []*Account `json:"-"`
}

type Transaction struct {
	Id          GUID
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
	Id             GUID
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
	return []byte(string('"') + (*big.Rat)(amt).RatString() + string('"')), nil
}

func (amt *Amount) UnmarshalJSON(s []byte) error {
	str, err := strconv.Unquote(string(s))
	if err != nil {
		return fmt.Errorf("invalid price string %q: %s", s, err)
	}
	_, ok := (*big.Rat)(amt).SetString(str)
	if !ok {
		return fmt.Errorf("invalid price string %q", str)
	}
	return nil
}

func (x *Amount) Add(y *Amount) *Amount { (*big.Rat)(x).Add((*big.Rat)(x), (*big.Rat)(y)); return x }

// An accounting book.
type Book struct {
	Accounts     map[GUID]*Account
	Transactions map[GUID]*Transaction

	// Computed data.
	Balance map[*Account]*Amount `json:"-"`
	// Flows by account.
	Flows map[*Account][]*Flow `json:"-"`
}

func (book *Book) Recompute() {
	book.Flows = book.sortFlows()
	book.Balance = make(map[*Account]*Amount, len(book.Accounts))
	for _, act := range book.Accounts {
		book.Balance[act] = sumFlows(book.Flows[act])
	}
}

func sumFlows(flows []*Flow) *Amount {
	total := new(Amount)
	for _, f := range flows {
		total = total.Add(f.Price)
	}
	return total
}

func (book *Book) sortFlows() (flows map[*Account][]*Flow) {
	flows = make(map[*Account][]*Flow, len(book.Accounts))
	for it, trn := range book.Transactions {
		for i, f := range trn.Flows {
			flows[f.Account] = append(flows[f.Account], &trn.Flows[i])
			trn.Flows[i].Parent = book.Transactions[it]
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

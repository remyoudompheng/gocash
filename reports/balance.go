package reports

import (
	"math/big"
	"sort"
	"time"

	"github.com/remyoudompheng/gocash/types"
)

// Balance produces a balance report out of a list of flows.
func Balance(flows [][]*types.Flow) BalanceReport {
	perMonth := make(map[int]*big.Rat, 20)

	for _, fs := range flows {
		for _, f := range fs {
			when := f.Parent.Date
			key := when.Year()*16 + int(when.Month())
			x := perMonth[key]
			if x == nil {
				x = new(big.Rat)
			}
			perMonth[key] = x.Add(x, (*big.Rat)(f.Price)) // FIXME: currency.
		}
	}

	var months []int
	for m := range perMonth {
		months = append(months, m)
	}
	sort.Ints(months)

	var rep BalanceReport
	val := new(big.Rat)
	for _, m := range months {
		t := time.Date(m/16, time.Month(m%16), 1,
			0, 0, 0, 0, time.UTC)
		next := new(big.Rat).Add(val, perMonth[m])
		rep.T = append(rep.T, t)
		rep.Values = append(rep.Values, next)
		val = next
	}

	return rep
}

type BalanceReport struct {
	T      []time.Time
	Values []*big.Rat
}

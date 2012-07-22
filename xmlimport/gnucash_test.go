package xmlimport

import (
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	const testfile = "testdata/test.xml"
	f, err := ReadFile(testfile)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", b)
}

func TestImport(t *testing.T) {
	const testfile = "testdata/test.xml"
	book, err := ImportFile(testfile)
	if err != nil {
		t.Fatal(err)
	}
	js, err := json.MarshalIndent(book, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", js)
}

func TestImportReal(t *testing.T) {
	const testfile = "/home/remy/Documents/banque/compta.xml"
	book, err := ImportFile(testfile)
	if err != nil {
		t.Fatal(err)
	}
	// Consistency check.
	t.Logf("%d accounts, %d transactions", len(book.Accounts), len(book.Transactions))
	book.Recompute()
	trnCount := make(map[*types.Account]int)
	for _, trans := range book.Transactions {
		for _, flow := range trans.Flows {
			trnCount[flow.Account]++
		}
	}

	total := new(big.Rat)
	for act, count := range trnCount {
		t.Logf("account %q: %d transactions, %s", act.Name, count, book.Balance[act].FloatString(2))
		total = total.Add(total, book.Balance[act])
	}
	t.Logf("total: %.2f (should be 0.0)", total.FloatString(2))
}

func parsePrice(s string) float64 {
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

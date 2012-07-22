package xmlimport

import (
	"encoding/json"
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
	const testfile = "testdata/test.xml"
	book, err := ImportFile(testfile)
	if err != nil {
		t.Fatal(err)
	}
	// Consistency check.
	t.Logf("%d accounts, %d transactions", len(book.Accounts), len(book.Transactions))
	grandtotal := 0.0
	for _, acct := range book.Accounts {
		count := 0
		total := 0.0
		for _, trans := range book.Transactions {
			for _, flow := range trans.Flows {
				if flow.Account == acct {
					count++
					total += parsePrice(flow.Price)
				}
			}
		}
		t.Logf("account %q: %d transactions, %.2f", acct.Name, count, total)
		grandtotal += total
	}
	t.Logf("total: %.2f (should be 0.0)", grandtotal)
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

package types

import (
	"encoding/json"
	"math/big"
	"testing"
)

func TestAmountJSON(t *testing.T) {
	type item struct {
		Price *Amount
	}
	const expected = `{"Price":"367/100"}`
	p := new(big.Rat)
	p.SetString("367/100")
	x := &item{Price: (*Amount)(p)}

	s, err := json.Marshal(x)
	if err != nil {
		t.Fatal(err)
	}
	if string(s) != expected {
		t.Errorf("got %q, expected %q", string(s), expected)
	}

	x = new(item)
	err = json.Unmarshal([]byte(expected), x)
	if err != nil {
		t.Fatal(err)
	}
	if str := x.Price.String(); str != "3.67" {
		t.Errorf("got %s, expected %s", str, "3.67")
	}
}

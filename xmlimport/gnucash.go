// Package xmlimport implements reading of the Gnucash 2 XML format.
package xmlimport

import (
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func Import(r io.Reader) (f *File, err error) {
	f = new(File)
	dec := xml.NewDecoder(r)
	err = dec.Decode(f)
	return
}

func ImportFile(name string) (f *File, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return Import(file)
}

type GUID string

func (g GUID) Bytes() (b [16]byte, err error) {
	x, err := hex.DecodeString(string(g))
	if len(x) != 16 && err == nil {
		err = fmt.Errorf("invalid GUID of length %s", len(x))
	}
	copy(b[:], x)
	return
}

type File struct {
	XMLName xml.Name
	Book    Book `xml:"book"`
}

type Book struct {
	XMLName      xml.Name      `xml:"http://www.gnucash.org/XML/gnc book"`
	Id           string        `xml:"id"`
	Commos       []Commodity   `xml:"http://www.gnucash.org/XML/gnc commodity"`
	Accounts     []Account     `xml:"http://www.gnucash.org/XML/gnc account"`
	Transactions []Transaction `xml:"http://www.gnucash.org/XML/gnc transaction"`
}

type Commodity struct {
	XMLName xml.Name
	Space   string `xml:"http://www.gnucash.org/XML/cmdty space"`
	Id      string `xml:"http://www.gnucash.org/XML/cmdty id"`
}

type Account struct {
	XMLName xml.Name `xml:"http://www.gnucash.org/XML/gnc account"`
	Name    string   `xml:"name"`
	Id      string   `xml:"id"`
	Type    string   `xml:"type"`
	Slots   Slots    `xml:"slots>slot"`
}

type Transaction struct {
	XMLName     xml.Name  `xml:"http://www.gnucash.org/XML/gnc transaction"`
	Id          string    `xml:"id"`
	Currency    Commodity `xml:"currency"`
	Description string    `xml:"description"`
	Splits      []Split   `xml:"splits>split"`
	Slots       Slots     `xml:"slots>slot"`
}

type Split struct {
	XMLName       xml.Name
	Id            string     `xml:"id"`
	Account       string     `xml:"account"`
	Reconciled    string     `xml:"reconciled-state"`
	ReconcileDate *TimeStamp `xml:"reconcile-date"`
}

// A slot represents a data structure similar to a JSON object.
type Slot struct {
	Key   string    `xml:"key"`
	Value SlotValue `xml:"value"`
}

type SlotValue struct {
	Type   string `xml:"type,attr"`
	String string `xml:",chardata"`
	Date   string `xml:"gdate"`
	Values Slots  `xml:"slot"`
}

type Slots []Slot

func (s Slots) Map() (m map[string]interface{}) {
	for _, slot := range s {
		val := slot.Value
		var v interface{}
		switch val.Type {
		case "integer":
			n, err := strconv.Atoi(val.String)
			if err != nil {
				panic(err)
			}
			v = n
		case "string":
			v = val.String
		case "frame":
			v = val.Values.Map()
		case "gdate":
			t, err := time.Parse("2006-01-02", val.Date)
			if err != nil {
				panic(err)
			}
			v = t
		default:
			err := fmt.Errorf("unknown slot type: %s", val.Type)
			panic(err)
		}
		if m == nil {
			m = make(map[string]interface{})
		}
		m[slot.Key] = v
	}
	return m
}

// MarshalJSON implements a JSON representation of a Slots as a
// simplified map. It is intended for debugging purposes.
func (s Slots) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Map())
}

type TimeStamp struct {
	Date string `xml:"http://www.gnucash.org/XML/ts date"`
	Ns   int    `xml:"http://www.gnucash.org/XML/ts ns"`
}

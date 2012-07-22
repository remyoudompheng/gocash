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

	"github.com/remyoudompheng/gocash/types"
)

// Read reads a File from r.
func Read(r io.Reader) (f *File, err error) {
	f = new(File)
	dec := xml.NewDecoder(r)
	err = dec.Decode(f)
	return
}

// ReadFile reads the named XML file.
func ReadFile(name string) (f *File, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Read(file)
}

// Import converts a Gnucash XML file read from r and
// returns a parsed accounting book.
func Import(r io.Reader) (book *types.Book, err error) {
	f, err := Read(r)
	if err == nil {
		book, err = f.Import()
	}
	return
}

// ImportFile converts the named Gnucash XML file and
// returns a parsed accounting book.
func ImportFile(name string) (book *types.Book, err error) {
	f, err := ReadFile(name)
	if err == nil {
		book, err = f.Import()
	}
	return
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

func (file *File) Import() (book *types.Book, err error) {
	book = new(types.Book)
	// Parse accounts.
	accountsById := make(map[GUID]*types.Account, len(file.Book.Accounts))
	parents := make(map[GUID]GUID, len(file.Book.Accounts))
	for _, xmlacct := range file.Book.Accounts {
		act, err := xmlacct.Import()
		if err != nil {
			return nil, err
		}
		accountsById[xmlacct.Id] = &act
		parents[xmlacct.Id] = xmlacct.Parent
		book.Accounts = append(book.Accounts, &act)
	}

	// Resolve account hierarchy.
	actNames := make(map[GUID]string)
	for _, xmlacct := range file.Book.Accounts {
		name := xmlacct.Name
		for t := xmlacct.Parent; t != ""; t = parents[t] {
			name = accountsById[t].Name + "/" + name
		}
		actNames[xmlacct.Id] = name
	}
	for guid, name := range actNames {
		accountsById[guid].Name = name
	}

	// Parse transactions.
	for _, xmltrn := range file.Book.Transactions {
		trn, err := xmltrn.Import(accountsById)
		if err != nil {
			return nil, fmt.Errorf("error in transaction %s: %s",
				xmltrn.Id, err)
		}
		book.Transactions = append(book.Transactions, trn)
	}
	return book, nil
}

type Book struct {
	XMLName      xml.Name      `xml:"http://www.gnucash.org/XML/gnc book"`
	Id           GUID          `xml:"id"`
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
	Id      GUID     `xml:"id"`
	Type    string   `xml:"type"`
	Slots   Slots    `xml:"slots>slot"`
	Parent  GUID     `xml:"parent"`
}

func (xmlact *Account) Import() (act types.Account, err error) {
	act = types.Account{
		Name: xmlact.Name,
		Type: xmlact.Type,
	}
	slots := xmlact.Slots.Map()
	if slots != nil && slots["notes"] != nil {
		if notes, ok := slots["notes"].(string); ok {
			act.Description = notes
		} else {
			return act, fmt.Errorf("description of account %s is not a string", act.Name)
		}
	}
	return act, nil
}

type Transaction struct {
	XMLName     xml.Name  `xml:"http://www.gnucash.org/XML/gnc transaction"`
	Id          string    `xml:"id"`
	Currency    Commodity `xml:"currency"`
	Number      string    `xml:"num"`
	Description string    `xml:"description"`
	Splits      []Split   `xml:"splits>split"`
	Slots       Slots     `xml:"slots>slot"`
	PostedDate  TimeStamp `xml:"date-posted"`
	EnteredDate TimeStamp `xml:"date-entered"`
}

func (xmltrn *Transaction) Import(accts map[GUID]*types.Account) (trn types.Transaction, err error) {
	trn.Description = xmltrn.Description
	trn.Number = xmltrn.Number
	if slots := xmltrn.Slots.Map(); slots != nil && slots["notes"] != nil {
		if notes, ok := slots["notes"].(string); ok {
			trn.Notes = notes
		} else {
			return trn, fmt.Errorf("description is not a string")
		}
	}

	trn.Flows = make([]types.Flow, len(xmltrn.Splits))
	for i, split := range xmltrn.Splits {
		trn.Flows[i], err = split.Import(accts)
		if err != nil {
			return trn, fmt.Errorf("error in split transaction %s: %s",
				split.Id, err)
		}
	}
	trn.Date, err = xmltrn.PostedDate.Time()
	if err != nil {
		return trn, fmt.Errorf("invalid date-posted: %s", err)
	}
	trn.Stamp, err = xmltrn.EnteredDate.Time()
	if err != nil {
		return trn, fmt.Errorf("invalid date-entered: %s", err)
	}
	return trn, nil
}

type Split struct {
	XMLName       xml.Name
	Id            GUID       `xml:"id"`
	Account       GUID       `xml:"account"`
	Memo          string     `xml:"memo"`
	Reconciled    string     `xml:"reconciled-state"`
	ReconcileDate *TimeStamp `xml:"reconcile-date"`
	Value         string     `xml:"value"`
	Quantity      string     `xml:"quantity"`
}

func (split *Split) Import(accts map[GUID]*types.Account) (flow types.Flow, err error) {
	flow.Account = accts[split.Account]
	if flow.Account == nil {
		return flow, fmt.Errorf("account %s does not exist", split.Account)
	}
	flow.Price = split.Value
	flow.Memo = split.Memo
	switch split.Reconciled {
	case "y":
		flow.Reconciled = true
		if split.ReconcileDate == nil {
			break
		}
		flow.ReconciledTime, err = split.ReconcileDate.Time()
		if err != nil {
			return flow, fmt.Errorf("invalid reconcile-date: %s", err)
		}
	case "n":
		flow.Reconciled = false
	default:
		return flow, fmt.Errorf("invalid reconciled state %q", split.Reconciled)
	}
	return flow, nil
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

func (s TimeStamp) Time() (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04:05 -0700", s.Date)
	return t.Add(time.Duration(s.Ns) * time.Nanosecond), err
}

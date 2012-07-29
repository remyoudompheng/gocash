// Package xmlimport implements reading of the Gnucash 2 XML format.
package xmlimport

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math/big"
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

type File struct {
	XMLName xml.Name
	Book    Book `xml:"book"`
}

func (file *File) Import() (book *types.Book, err error) {
	book = new(types.Book)
	// Parse accounts.
	accountsById := make(map[types.GUID]*types.Account, len(file.Book.Accounts))
	parents := make(map[types.GUID]types.GUID, len(file.Book.Accounts))
	for _, xmlacct := range file.Book.Accounts {
		act, err := xmlacct.Import()
		if err != nil {
			return nil, err
		}
		accountsById[xmlacct.Id] = &act
		parents[xmlacct.Id] = xmlacct.Parent
	}
	book.Accounts = accountsById

	// Resolve account hierarchy.
	actNames := make(map[types.GUID]string)
	for _, xmlacct := range file.Book.Accounts {
		name := xmlacct.Name
		for t := xmlacct.Parent; t != ""; t = parents[t] {
			act := accountsById[t]
			if act.Type == "ROOT" {
				name = "/" + name
			} else {
				name = accountsById[t].Name + "/" + name
			}
		}
		actNames[xmlacct.Id] = name
		if parent := accountsById[xmlacct.Parent]; parent != nil {
			parent.Children = append(parent.Children, accountsById[xmlacct.Id])
		}
	}
	for guid, name := range actNames {
		accountsById[guid].Name = name
	}

	// Parse transactions.
	book.Transactions = make(map[types.GUID]*types.Transaction, len(file.Book.Transactions))
	for _, xmltrn := range file.Book.Transactions {
		trn, err := xmltrn.Import(accountsById)
		if err != nil {
			return nil, fmt.Errorf("error in transaction %s: %s",
				xmltrn.Id, err)
		}
		book.Transactions[xmltrn.Id] = trn
	}

	return book, nil
}

type Book struct {
	XMLName      xml.Name      `xml:"http://www.gnucash.org/XML/gnc book"`
	Id           types.GUID    `xml:"id"`
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
	XMLName   xml.Name   `xml:"http://www.gnucash.org/XML/gnc account"`
	Name      string     `xml:"name"`
	Id        types.GUID `xml:"id"`
	Commodity Commodity  `xml:"commodity"`
	Type      string     `xml:"type"`
	Slots     Slots      `xml:"slots>slot"`
	Parent    types.GUID `xml:"parent"`
}

func (xmlact *Account) Import() (act types.Account, err error) {
	act = types.Account{
		Id:   xmlact.Id,
		Name: xmlact.Name,
		Type: xmlact.Type,
		Unit: xmlact.Commodity.Id,
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
	XMLName     xml.Name   `xml:"http://www.gnucash.org/XML/gnc transaction"`
	Id          types.GUID `xml:"id"`
	Currency    Commodity  `xml:"currency"`
	Number      string     `xml:"num"`
	Description string     `xml:"description"`
	Splits      []Split    `xml:"splits>split"`
	Slots       Slots      `xml:"slots>slot"`
	PostedDate  TimeStamp  `xml:"date-posted"`
	EnteredDate TimeStamp  `xml:"date-entered"`
}

func (xmltrn *Transaction) Import(accts map[types.GUID]*types.Account) (trn *types.Transaction, err error) {
	trn = &types.Transaction{
		Id:          xmltrn.Id,
		Description: xmltrn.Description,
		Number:      xmltrn.Number,
	}
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
	Id            types.GUID `xml:"id"`
	Account       types.GUID `xml:"account"`
	Memo          string     `xml:"memo"`
	Reconciled    string     `xml:"reconciled-state"`
	ReconcileDate *TimeStamp `xml:"reconcile-date"`
	Value         string     `xml:"value"`
	Quantity      string     `xml:"quantity"`
}

func (split *Split) Import(accts map[types.GUID]*types.Account) (flow types.Flow, err error) {
	flow = types.Flow{
		Id:      split.Id,
		Account: accts[split.Account],
		Price:   new(types.Amount),
		Memo:    split.Memo,
	}

	if flow.Account == nil {
		return flow, fmt.Errorf("account %s does not exist", split.Account)
	}
	_, ok := (*big.Rat)(flow.Price).SetString(split.Value)
	if !ok {
		return flow, fmt.Errorf("incorrect price format: %q", split.Value)
	}
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

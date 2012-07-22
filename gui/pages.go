package gui

import (
	"fmt"
	"io"
	"net/http"

	"github.com/remyoudompheng/gocash/types"
)

func pageHome(book *types.Book, w io.Writer, req *http.Request) error {
	return homeTpl.Execute(w, templateData{
		Title: "Gocash",
		Book:  book,
	})
}

func pageBook(book *types.Book, w io.Writer, req *http.Request) error {
	return bookTpl.Execute(w, templateData{
		Title: "Book",
		Book:  book,
	})
}

func pageAccount(book *types.Book, w io.Writer, req *http.Request) error {
	req.ParseForm()
	acctname := req.Form.Get("name")
	var account *types.Account
	for _, acct := range book.Accounts {
		if acct.Name == acctname {
			account = acct
		}
	}
	if account == nil {
		return fmt.Errorf("no such account: %q", acctname)
	}

	return accountTpl.Execute(w, templateData{
		Title:   "Account",
		Book:    book,
		Account: account,
	})
}

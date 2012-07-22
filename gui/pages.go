package gui

import (
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
	return accountTpl.Execute(w, templateData{
		Title: "Account",
		Book:  book,
	})
}

package gui

import (
	"bytes"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"path/filepath"

	"github.com/remyoudompheng/go-misc/weblibs"

	"github.com/remyoudompheng/gocash/types"
	"io"
)

func StartServer(addr string, book *types.Book) error {
	parseTemplates()
	err := weblibs.RegisterAll(http.DefaultServeMux)
	if err != nil {
		return err
	}
	http.Handle("/", curryBook(book, pageHome))
	http.Handle("/account/", curryBook(book, pageAccount))
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(StaticDir)),
	))
	log.Printf("starting HTTP server at %s", addr)
	return http.ListenAndServe(addr, nil)
}

type httpHandler func(*types.Book, io.Writer, *http.Request) error

// curryBook makes http handlers out of handlers parameterized
// by an accounting book.
func curryBook(book *types.Book, h httpHandler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			log.Printf("%s %s from %s", req.Method, req.URL, req.RemoteAddr)
			resp := new(bytes.Buffer)
			err := h(book, resp, req)
			if err == nil {
				w.Write(resp.Bytes())
			} else {
				log.Printf("ERROR: %s", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
}

var StaticDir = "static/"

func tplPath(name string) string { return filepath.Join(StaticDir, "templates", name+".tpl") }

func parseTemplate(name string) (*template.Template, error) {
	return template.New(name).
		Funcs(template.FuncMap{"cumul": cumulFlows}).
		ParseFiles(tplPath("common"), tplPath(name))
}

func cumulFlows(flows []*types.Flow) (bals []*types.Amount) {
	x := new(big.Rat)
	for _, f := range flows {
		x = x.Add(x, (*big.Rat)(f.Price))
		bals = append(bals, new(types.Amount).SetRat(x))
	}
	return
}

var homeTpl, bookTpl, accountTpl *template.Template

func parseTemplates() {
	homeTpl = template.Must(parseTemplate("home")).Lookup("common")
	bookTpl = template.Must(parseTemplate("book")).Lookup("common")
	accountTpl = template.Must(parseTemplate("account")).Lookup("common")
}

type templateData struct {
	Title   string
	Book    *types.Book
	Account *types.Account
}

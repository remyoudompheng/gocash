// gocash is a personal accounting application similar to gnucash
package main

import (
	"flag"
	"log"

	"github.com/remyoudompheng/gocash/gui"
	"github.com/remyoudompheng/gocash/xmlimport"
)

func init() {
	log.SetPrefix("gocash ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	var (
		filename string
		httpAddr string
	)
	flag.StringVar(&filename, "f", "", "path to GNucash XML file")
	flag.StringVar(&httpAddr, "http", "localhost:8099", "address of HTTP server")
	flag.StringVar(&gui.StaticDir, "static", "static/", "path to static files")
	flag.Parse()

	book, err := xmlimport.ImportFile(filename)
	if err != nil {
		log.Fatalf("ERROR: failed to load %q: %s", filename,  err)
	}
	log.Printf("Loaded %q: %d accounts, %d transactions",
		filename, len(book.Accounts), len(book.Transactions))

	err = gui.StartServer(httpAddr, book)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
}

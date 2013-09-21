// gocash is a personal accounting application similar to gnucash
package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/remyoudompheng/gocash/gui"
	"github.com/remyoudompheng/gocash/reports"
	"github.com/remyoudompheng/gocash/types"
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

		report string
	)
	flag.StringVar(&filename, "f", "", "path to GNucash XML file")
	flag.StringVar(&httpAddr, "http", "localhost:8099", "address of HTTP server")
	flag.StringVar(&gui.StaticDir, "static", "static/", "path to static files")
	flag.StringVar(&report, "report", "", "make a report")
	flag.Parse()

	t0 := time.Now()
	book, err := xmlimport.ImportFile(filename)
	if err != nil {
		log.Fatalf("ERROR: failed to load %q: %s", filename, err)
	}
	log.Printf("Loaded %q: %d accounts, %d transactions, in %s",
		filename, len(book.Accounts), len(book.Transactions), time.Since(t0))
	book.Recompute()

	switch {
	case report != "":
		// make a report.
		switch report {
		case "totalassets":
			var assetFlows [][]*types.Flow
			for acc, flows := range book.Flows {
				switch acc.Type {
				case "BANK", "ASSET", "CASH":
					assetFlows = append(assetFlows, flows)
				}
			}
			r := reports.Balance(assetFlows)
			for i := range r.T {
				fmt.Printf("%s,%s\n", r.T[i].Format("Jan 2006"), (*types.Amount)(r.Values[i]))
			}
		default:
			flag.Usage()
		}
	case httpAddr != "":
		err = gui.StartServer(httpAddr, book)
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
	default:
		flag.Usage()
	}
}

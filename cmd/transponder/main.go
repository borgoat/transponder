package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"

	"github.com/transponder-tf/transponder/pkg/server"
	"github.com/transponder-tf/transponder/pkg/statemgrmap"
)

var (
	flagData string

	flagAddress string
)

func init() {
	flag.StringVar(&flagData, "data", os.TempDir(), "path where the db file will be stored")

	flag.StringVar(&flagAddress, "address", ":1492", "address:port to bind the listener")
}

func main() {

	if !flag.Parsed() {
		flag.Parse()
	}

	srv := &http.Server{
		Addr:         flagAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	statePath := path.Join(flagData, ".state")
	mgrmap := statemgrmap.NewFilesystemMapWithPrefix(statePath)
	log.Printf("State stored in %s", statePath)

	bs := server.NewHTTPBackendServer(mgrmap)

	r := mux.NewRouter()
	srv.Handler = r

	bs.HandleWithRouter(r.PathPrefix("/terraform").Subrouter())

	log.Printf("Starting server on %s", srv.Addr)
	srv.ListenAndServe()
}
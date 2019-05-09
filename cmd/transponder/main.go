package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/mux"

	"github.com/transponder-tf/transponder/pkg/states/statemgrmap"

	"github.com/philippgille/gokv/file"
)

var (
	flagData string

	flagAddress string
)

func init() {
	flag.StringVar(
		&flagData,
		"data",
		path.Join(os.TempDir(), "transponder"),
		"path where the db file will be stored")

	flag.StringVar(
		&flagAddress,
		"address",
		":1492",
		"address:port to bind the listener")
}

func main() {

	if !flag.Parsed() {
		flag.Parse()
	}

	options := file.DefaultOptions
	options.Directory = flagData

	log.Printf("[DEBUG] Creating gokv.Store with %T %+v", options, options)
	store, err := file.NewStore(options)

	if err != nil {
		log.Fatalf("[ERROR] Could not create gokv.Store: %v", err)
	}

	log.Printf("[DEBUG] Created gokv.Store with type %T", store)

	mgrmap := statemgrmap.NewGoKVMap(store)

	bs := newHTTPBackendServer(mgrmap)

	
	srv := &http.Server{
		Addr:         flagAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	r := mux.NewRouter()
	srv.Handler = r
	log.Printf("[DEBUG] Created http.Server with options %+v", srv)

	log.Printf("[DEBUG] Handling /terraform with %T", bs)
	bs.HandleWithRouter(r.PathPrefix("/terraform").Subrouter())

	log.Printf("[INFO] Starting server on %s", srv.Addr)
	srv.ListenAndServe()
}

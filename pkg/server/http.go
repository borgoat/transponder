package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/states/statefile"
	"github.com/hashicorp/terraform/states/statemgr"

	"github.com/transponder-tf/transponder/pkg/statemgrmap"
)

// HttpBackendServer is an implementation compatible with
// https://www.terraform.io/docs/backends/types/http.html
type HttpBackendServer struct {
	mgrmap *statemgrmap.StateMgrMap
}

func NewHttpBackendServer(stateMgrMap *statemgrmap.StateMgrMap) *HttpBackendServer {

	return &HttpBackendServer{
		mgrmap: stateMgrMap,
	}
}

func (bs *HttpBackendServer) GetHandler() http.Handler {
	s := mux.NewRouter().Path("/{namespace}").Subrouter()

	s.Methods("GET").HandlerFunc(bs.getHandler)
	s.Methods("POST").HandlerFunc(bs.postHandler)
	s.Methods("DELETE").HandlerFunc(bs.deleteHandler)

	s.Methods("LOCK").HandlerFunc(bs.lockHandler)
	s.Methods("UNLOCK").HandlerFunc(bs.unlockHandler)

	return s
}

func (bs *HttpBackendServer) HandleWithRouter(r *mux.Router) {
	s := r.Path("/{namespace}").Subrouter()

	s.Methods("GET").HandlerFunc(bs.getHandler)
	s.Methods("POST").HandlerFunc(bs.postHandler)
	s.Methods("DELETE").HandlerFunc(bs.deleteHandler)

	s.Methods("LOCK").HandlerFunc(bs.lockHandler)
	s.Methods("UNLOCK").HandlerFunc(bs.unlockHandler)
}

func (bs *HttpBackendServer) stateMgrFromRequest(req *http.Request) statemgr.Full {
	vars := mux.Vars(req)

	mgr, err := bs.mgrmap.Get(vars["namespace"])
	if err != nil {
		panic(err)
	}

	return mgr
}

func (bs *HttpBackendServer) getHandler(res http.ResponseWriter, req *http.Request) {
	mgr := bs.stateMgrFromRequest(req)

	// We ensure our state manager has the latest version
	mgr.RefreshState()

	// We now export the statefile from our state manager
	sf := statemgr.Export(mgr)

	// Marshal it to a buffer and copy it to our HTTP response
	var buf bytes.Buffer

	// Only attempt to write the statefile to the buffer if not empty
	if sf != nil {
		statefile.Write(sf, &buf)
	}

	io.Copy(res, &buf)
}

func (bs *HttpBackendServer) postHandler(res http.ResponseWriter, req *http.Request) {
	mgr := bs.stateMgrFromRequest(req)

	// We read the statefile as it comes in the request body
	sf, err := statefile.Read(req.Body)
	if err != nil {
		panic(err)
	}

	// We now import it via statemgr to our backend state manager
	statemgr.Import(sf, mgr, false)

	// and ensure the state is persisted
	mgr.PersistState()
}

func (bs *HttpBackendServer) deleteHandler(res http.ResponseWriter, req *http.Request) {
	mgr := bs.stateMgrFromRequest(req)

	statemgr.Import(statemgr.NewStateFile(), mgr, true)
	mgr.PersistState()
}

func (bs *HttpBackendServer) lockHandler(res http.ResponseWriter, req *http.Request) {
	mgr := bs.stateMgrFromRequest(req)

	info, err := lockInfoFromRequest(req)
	if err != nil {
		// TODO Handle case where lock info is not valid
		panic(err)
	}

	lockID, err := mgr.Lock(info)
	if err != nil {
		// TODO Get LockError and answer correctly
		panic(err)
	}
	log.Printf("LockID: %s", lockID)
}

func (bs *HttpBackendServer) unlockHandler(res http.ResponseWriter, req *http.Request) {
	mgr := bs.stateMgrFromRequest(req)

	info, err := lockInfoFromRequest(req)
	if err != nil {
		panic(err)
	}

	mgr.Unlock(info.ID)
}

func lockInfoFromRequest(req *http.Request) (*statemgr.LockInfo, error) {
	// Read the LockInfo from the body
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	// and unmarshal it into the correct type
	info := &statemgr.LockInfo{}
	err := json.Unmarshal(buf.Bytes(), info)

	return info, err
}

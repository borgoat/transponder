package statemgr

import (
	"strings"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform/states/statemgr"
	"github.com/philippgille/gokv"
)

// GoKV implements statemgr.Full
// and stores state in a backend
// defined as a gokv.Store
type GoKV struct {
	kv gokv.Store

	path string

	current *states.State
}

var _ statemgr.Full = (*GoKV)(nil)

// NewGoKV creates a statemgr.Full based on philippgille/gokv
func NewGoKV(store gokv.Store, namespace string) *GoKV {
	return &GoKV{
		kv: store,
		path: namespace,
	}
}

func (s *GoKV) Lock(info *statemgr.LockInfo) (string, error) {
	info.ID = strings.Join([]string{s.path, info.ID}, "/")
	s.kv.Set(info.ID, info)

	return info.ID, nil
}

func (s *GoKV) Unlock(id string) error {
	s.kv.Delete(id)

	return nil
}

func (s *GoKV) State() *states.State {
	state := states.NewState()

	if (s.current != nil) {
		*state = *s.current
	}

	return state
}

func (s *GoKV) WriteState(state *states.State) error {
	s.current = state

	return nil
}

func (s *GoKV) RefreshState() error {
	s.kv.Get(s.path, s.current)

	return nil
}

func (s *GoKV) PersistState() error {
	s.kv.Set(s.path, s.current)

	return nil
}
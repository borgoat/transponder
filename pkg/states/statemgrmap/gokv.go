package statemgrmap

import (
	"github.com/hashicorp/terraform/states/statemgr"

	"github.com/philippgille/gokv"
	
	kvStatemgr "github.com/transponder-tf/transponder/pkg/states/statemgr"
)

type goKVStateLoader struct {
	Store gokv.Store
}

func (sl *goKVStateLoader) newGoKVStateMgr(namespace string) statemgr.Full {
	return kvStatemgr.NewGoKV(sl.Store, namespace)
}

func NewGoKVMap(store gokv.Store) *StateMgrMap {
	sl := &goKVStateLoader{
		Store: store,
	}

	return &StateMgrMap{
		newWithNamespace: sl.newGoKVStateMgr,
	}
}
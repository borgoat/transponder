package statemgrmap

import (
	"errors"

	"github.com/hashicorp/terraform/states/statemgr"
)

// NewMgrWithNamespace is a function returning a statemgr.Full
// for each "namespace" given, as a simple way to achieve multitenancy
type NewMgrWithNamespace func(namespace string) statemgr.Full

// StateMgrMap is an object capable of creating and providing
// multiple statemgr.Full instances, based on a key (or "namespace")
type StateMgrMap struct {
	mgrs             map[string]statemgr.Full
	newWithNamespace NewMgrWithNamespace
}

// Get (instantiate if needed) the statemgr.Full for the given namespace
func (mgrmap *StateMgrMap) Get(namespace string) (statemgr.Full, error) {
	if mgrmap.mgrs == nil {
		mgrmap.mgrs = make(map[string]statemgr.Full)
	}

	if namespace == "" {
		return nil, errors.New("namespace cannot be empty")
	}

	if mgrmap.mgrs[namespace] == nil {
		mgrmap.mgrs[namespace] = mgrmap.newWithNamespace(namespace)
	}

	return mgrmap.mgrs[namespace], nil
}

package statemgrmap

import (
	"errors"

	"github.com/hashicorp/terraform/states/statemgr"
)

type NewMgrWithNamespace func(namespace string) statemgr.Full

type StateMgrMap struct {
	mgrs             map[string]statemgr.Full
	newWithNamespace NewMgrWithNamespace
}

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

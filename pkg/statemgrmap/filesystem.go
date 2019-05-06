package statemgrmap

import (
	"path"

	"github.com/hashicorp/terraform/states/statemgr"
)

type filesystemStateLoader struct {
	Prefix string
}

func (fsl filesystemStateLoader) newFilesystemStateMgr(namespace string) statemgr.Full {
	return statemgr.NewFilesystem(path.Join(fsl.Prefix, namespace + ".tfstate"))
}

// NewFilesystemMap creates a StateMgrMap using statemgr.Filesystem as a state manager
func NewFilesystemMap(prefix string) *StateMgrMap {
	sl := &filesystemStateLoader{
		Prefix: prefix,
	}

	mgrmap := &StateMgrMap{
		newWithNamespace: sl.newFilesystemStateMgr,
	}

	return mgrmap
}

package statemgrmap

import (
	"github.com/hashicorp/terraform/states/statemgr"
)

type filesystemStateLoader struct {
	Prefix string
}

func (fsl filesystemStateLoader) NewFilesystemWithNamespace(namespace string) statemgr.Full {
	return statemgr.NewFilesystem(fsl.Prefix + "/" + namespace + ".tfstate")
}

func NewFilesystemMapWithPrefix(prefix string) *StateMgrMap {
	sl := &filesystemStateLoader{
		Prefix: prefix,
	}

	mgrmap := &StateMgrMap{
		newWithNamespace: sl.NewFilesystemWithNamespace,
	}

	return mgrmap
}

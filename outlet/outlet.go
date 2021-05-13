package outlet

import (
	"errors"

	"github.com/funcas/cgs/manager"
)

type Outlet struct {
	name             string
	acceptTransCodes []string
	executor         Executor
}

func (o Outlet) Executor() Executor {
	return o.executor
}

func (o Outlet) Name() string {
	return o.name
}

func (o Outlet) AcceptTransCodes() []string {
	return o.acceptTransCodes
}

func NewOutlet(name string, executor Executor, acceptTransCodes []string) *Outlet {
	return &Outlet{
		name:             name,
		executor:         executor,
		acceptTransCodes: acceptTransCodes,
	}
}

type TransCodeMap map[string]string

func (t TransCodeMap) Exists(transCode string) bool {
	_, exists := t[transCode]
	return exists
}

func (t TransCodeMap) GetOutlet(transCode string) (*Outlet, error) {
	if !t.Exists(transCode) {
		return nil, errors.New("invalid transCode")
	}

	res, err := manager.GetDiFactory().SafeGet(t[transCode])
	if err != nil {
		return nil, err
	}
	return res.(*Outlet), nil
}

package outlet

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

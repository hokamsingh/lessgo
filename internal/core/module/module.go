package module

type Module struct {
	Name        string
	Controllers []interface{}
	Services    []interface{}
}

func NewModule(name string, controllers []interface{}, services []interface{}) *Module {
	return &Module{
		Name:        name,
		Controllers: controllers,
		Services:    services,
	}
}

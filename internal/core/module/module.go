package module

type IModule interface {
	GetName() string
	GetControllers() []interface{}
	GetServices() []interface{}
}

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

func (m *Module) GetName() string {
	return m.Name
}

func (m *Module) GetControllers() []interface{} {
	return m.Controllers
}

func (m *Module) GetServices() []interface{} {
	return m.Services
}

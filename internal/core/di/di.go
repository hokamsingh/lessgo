package di

import (
	"go.uber.org/dig"
)

type Container struct {
	container *dig.Container
}

func NewContainer() *Container {
	return &Container{
		container: dig.New(),
	}
}

func (c *Container) Register(constructor interface{}) error {
	return c.container.Provide(constructor)
}

// Provide is an alias for Register. It registers a constructor or provider in the container
func (c *Container) Provide(constructor interface{}) error {
	return c.Register(constructor)
}

func (c *Container) Invoke(function interface{}) error {
	return c.container.Invoke(function)
}

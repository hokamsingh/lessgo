package web

import (
	"errors"
	"reflect"
	"sync"
)

type Container struct {
	mu        sync.RWMutex
	instances map[string]reflect.Value
	providers map[string]reflect.Value
}

func NewContainer() *Container {
	return &Container{
		instances: make(map[string]reflect.Value),
		providers: make(map[string]reflect.Value),
	}
}

func (c *Container) Register(name string, provider interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.providers[name]; exists {
		return errors.New("provider already registered")
	}

	c.providers[name] = reflect.ValueOf(provider)
	return nil
}

func (c *Container) Resolve(name string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if instance, exists := c.instances[name]; exists {
		return instance.Interface(), nil
	}

	provider, exists := c.providers[name]
	if !exists {
		return nil, errors.New("provider not found")
	}

	// Call the provider to get a new instance
	result := provider.Call([]reflect.Value{})
	if len(result) != 1 {
		return nil, errors.New("provider must return exactly one value")
	}

	instance := result[0]
	c.instances[name] = instance
	return instance.Interface(), nil
}

func (c *Container) RegisterMultiple(name string, instances ...interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, instance := range instances {
		instanceName := name + "_" + reflect.TypeOf(instance).Name()
		if _, exists := c.instances[instanceName]; exists {
			return errors.New("instance already registered")
		}
		c.instances[instanceName] = reflect.ValueOf(instance)
	}
	return nil
}

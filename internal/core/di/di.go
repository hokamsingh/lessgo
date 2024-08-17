package di

import (
	"errors"
	"reflect"
	"sync"
)

type Container struct {
	services map[string]interface{}
	mu       sync.RWMutex
}

func NewContainer() *Container {
	return &Container{
		services: make(map[string]interface{}),
	}
}

func (c *Container) Register(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

func (c *Container) Get(name string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	service, exists := c.services[name]
	if !exists {
		return nil, errors.New("service not found")
	}
	return service, nil
}

func (c *Container) Inject(target interface{}) error {
	val := reflect.ValueOf(target).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.CanSet() && field.Kind() == reflect.Interface {
			serviceName := field.Type().String()
			service, err := c.Get(serviceName)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(service))
		}
	}
	return nil
}

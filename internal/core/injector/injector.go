package injector

import (
	"errors"
	"reflect"
	"sync"
)

type Scope int

const (
	Singleton Scope = iota
	Transient
	Scoped
)

type Provider struct {
	Factory reflect.Value
	Scope   Scope
}

type Container struct {
	mu        sync.RWMutex
	providers map[string]Provider
	instances map[string]reflect.Value
	scoped    map[string]map[string]reflect.Value // Scoped instances are tracked per provider name
}

func NewContainer() *Container {
	return &Container{
		providers: make(map[string]Provider),
		instances: make(map[string]reflect.Value),
		scoped:    make(map[string]map[string]reflect.Value),
	}
}

// Register registers a provider in the container with a specific name and scope.
func (c *Container) Register(name string, provider interface{}, scope Scope) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.providers[name]; exists {
		return errors.New("provider already registered")
	}

	c.providers[name] = Provider{
		Factory: reflect.ValueOf(provider),
		Scope:   scope,
	}

	// Initialize scoped map for the provider
	if scope == Scoped {
		c.scoped[name] = make(map[string]reflect.Value)
	}

	return nil
}

// Resolve resolves an instance by its name and handles different scopes.
func (c *Container) Resolve(name string, scopeID ...string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	provider, exists := c.providers[name]
	if !exists {
		return nil, errors.New("provider not found")
	}

	// Handle singleton scope
	if provider.Scope == Singleton {
		if instance, exists := c.instances[name]; exists {
			return instance.Interface(), nil
		}
	}

	// Handle scoped scope
	if provider.Scope == Scoped {
		if len(scopeID) == 0 {
			return nil, errors.New("scope ID is required for scoped providers")
		}
		if instance, exists := c.scoped[name][scopeID[0]]; exists {
			return instance.Interface(), nil
		}
	}

	// Create new instance
	result := provider.Factory.Call([]reflect.Value{})
	if len(result) != 1 {
		return nil, errors.New("provider must return exactly one value")
	}
	instance := result[0]

	// Store the instance according to the scope
	if provider.Scope == Singleton {
		c.instances[name] = instance
	} else if provider.Scope == Scoped {
		c.scoped[name][scopeID[0]] = instance
	}

	return instance.Interface(), nil
}

// Cleanup cleans up all instances in the container.
func (c *Container) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for name, instance := range c.instances {
		if lifecycle, ok := instance.Interface().(Lifecycle); ok {
			_ = lifecycle.Cleanup()
		}
		delete(c.instances, name)
	}

	for name, scopedInstances := range c.scoped {
		for scopeID, instance := range scopedInstances {
			if lifecycle, ok := instance.Interface().(Lifecycle); ok {
				_ = lifecycle.Cleanup()
			}
			delete(c.scoped[name], scopeID)
		}
	}
}

// Lifecycle interface provides hooks for initialization and cleanup.
type Lifecycle interface {
	Init() error
	Cleanup() error
}

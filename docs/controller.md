
---

## Package `controller`

The `controller` package provides a base structure and interface for defining and registering routes in the application.

### Types

#### `Controller`

The `Controller` interface defines the contract for all controllers in the application. Controllers implementing this interface must define the `RegisterRoutes` method to set up the necessary routes.

```go
type Controller interface {
    RegisterRoutes(r *router.Router)
}
```

#### `BaseController`

The `BaseController` struct provides a default implementation of the `Controller` interface. It can be embedded in other controllers to inherit default behavior or overridden with custom implementations.

```go
type BaseController struct {
}
```

### Methods

#### `BaseController.RegisterRoutes`

```go
func (bc *BaseController) RegisterRoutes(r *router.Router)
```

The `RegisterRoutes` method provides a default implementation of the `Controller` interface's method. It can be overridden by embedding `BaseController` in another struct and defining a custom implementation.

**Example:**

```go
type TestController struct {
    BaseController
    Path    string
    Service TestService
}

func NewTestController(service *TestService, path string) *TestController {
    return &TestController{
        Service: *service,
        Path:    path,
    }
}

func (tc *TestController) RegisterRoutes(r *router.Router) {
    tr := r.SubRouter(tc.Path)
    tr.Get("/ping", func(ctx *router.Context) {
        ctx.Send("pong")
    })
}
```

#### `RegisterModuleRoutes`

```go
func RegisterModuleRoutes(r *router.Router, m module.IModule)
```

The `RegisterModuleRoutes` function registers routes for a module. It iterates through the module's controllers, ensuring each controller implements the `Controller` interface and then calls `RegisterRoutes` to set up the routes.

**Parameters:**

- `r`: The router instance used to register the routes.
- `m`: The module containing controllers to register.

**Usage:**

```go
RegisterModuleRoutes(routerInstance, moduleInstance)
```

---
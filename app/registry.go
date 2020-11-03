package app

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/golog"
    "github.com/kataras/iris/v12/mvc"
    "reflect"
    "sync"
)

type ComponentsRegistry struct {
    sync.RWMutex
    components []interface{}
}

func NewComponentsRegistry() *ComponentsRegistry {
    return &ComponentsRegistry{
        components: make([]interface{}, 0)}
}

func (r *ComponentsRegistry) Register(component interface{}) {
    r.Lock()
    defer r.Unlock()
    if nil == component {
        return
    }
    r.components = append(r.components, component)
}

var (
    dependencies = NewComponentsRegistry()
    middlewares  = NewComponentsRegistry()
    controllers  = NewComponentsRegistry()
)

func RegisterDependency(dependency interface{}) {
    dependencies.Register(dependency)
}

func RegisterMiddleware(middleware interface{}) {
    middlewares.Register(middleware)
}

func RegisterController(controller interface{}) {
    controllers.Register(controller)
}

type Components struct {
    config       *base.Config
    db           *sqlx.DB
    dependencies []interface{}
    middlewares  []interface{}
    controllers  []interface{}
}

func PrepareComponents(config *base.Config, db *sqlx.DB) *Components {
    components := &Components{config: config, db: db,
        dependencies: make([]interface{}, 0),
        middlewares:  make([]interface{}, 0),
        controllers:  make([]interface{}, 0)}
    components.dependencies = append(
        components.dependencies, config, db)

    for _, dependency := range dependencies.components {
        resolved := components.resolveDependencies(dependency)
        components.dependencies = append(components.dependencies, resolved...)
    }
    for _, middleware := range middlewares.components {
        components.middlewares = append(components.middlewares, middleware)
    }
    for _, controller := range controllers.components {
        components.controllers = append(components.controllers, controller)
    }

    return components
}

var (
    configType = reflect.TypeOf((*base.Config)(nil))
    dbType     = reflect.TypeOf((*sqlx.DB)(nil))
)

func (c *Components) resolveDependencies(dependency interface{}) []interface{} {
    v := reflect.ValueOf(dependency)
    if reflect.Func != v.Kind() {
        return []interface{}{dependency}
    }

    t := v.Type()
    numIn := t.NumIn()
    if numIn > 2 {
        golog.Error("Not Supported component builder func inputs")
        return []interface{}{}
    }

    inputs := make([]reflect.Value, numIn)
    for i := 0; i < numIn; i++ {
        inputType := t.In(i)
        if configType == inputType {
            inputs[i] = reflect.ValueOf(c.config)
        } else if dbType == inputType {
            inputs[i] = reflect.ValueOf(c.db)
        }
    }
    outputs := v.Call(inputs)
    resolved := make([]interface{}, len(outputs))
    for i, output := range outputs {
        resolved[i] = output.Interface()
    }
    return resolved
}

func (c *Components) Configurator(app *mvc.Application) {
    router := app.Router.ConfigureContainer()
    routerContainer := router.Container

    for _, dependency := range c.dependencies {
        // register twice because of app.container
        // is different from app.Router.Container
        app.Register(dependency)
        routerContainer.Register(dependency)
    }
    for _, middleware := range c.middlewares {
        router.Use(middleware)
    }
    for _, controller := range c.controllers {
        app.Handle(controller)
    }
}

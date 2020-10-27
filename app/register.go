package app

import (
    . "github.com/CharLemAznable/go-lannister/elf"
    "github.com/kataras/iris/v12/mvc"
)

var (
    dependencies = NewRegistry("Dependency")
    middlewares  = NewRegistry("Middleware")
    controllers  = NewRegistry("Controller")
)

func RegisterDependency(name string, dependency interface{}) {
    dependencies.Register(name, dependency)
}

func RegisterMiddleware(name string, middleware interface{}) {
    middlewares.Register(name, middleware)
}

func RegisterController(name string, controller interface{}) {
    controllers.Register(name, controller)
}

func dependencyConfigurator(app *mvc.Application) {
    router := app.Router.ConfigureContainer().Container
    dependencies.IterateSorted(func(dependency interface{}) {
        // register twice because of app.container
        // is different from app.Router.Container
        app.Register(dependency)
        router.Register(dependency)
    })
}

func middlewareConfigurator(app *mvc.Application) {
    router := app.Router.ConfigureContainer()
    middlewares.Iterate(func(middleware interface{}) {
        router.Use(middleware)
    })
}

func controllerConfigurator(app *mvc.Application) {
    controllers.Iterate(func(controller interface{}) {
        app.Handle(controller)
    })
}

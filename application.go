package lannister

import (
    "github.com/CharLemAznable/gokits"
    "github.com/kataras/golog"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/middleware/logger"
    "github.com/kataras/iris/v12/middleware/recover"
    "github.com/kataras/iris/v12/middleware/requestid"
    "github.com/kataras/iris/v12/mvc"
)

type application struct {
    app *iris.Application
}

func Application() *application {
    preparingDB()

    app := iris.New()
    app.Logger().SetLevel(appConfig.LogLevel)
    app.Configure(iris.WithSocketSharding,
        iris.WithoutBodyConsumptionOnUnmarshal)
    app.UseRouter(requestid.New())
    app.UseRouter(recover.New())
    if app.Logger().Level == golog.DebugLevel {
        app.UseRouter(logger.New())
    }
    app.Get("/", func(ctx iris.Context) {})
    mvc.Configure(app.Party(appConfig.ContextPath),
        dependencyConfigurator,
        middlewareConfigurator,
        controllerConfigurator)

    application := new(application)
    application.app = app
    return application
}

func (a *application) Run() {
    _ = a.app.Listen(":" + gokits.StrFromInt(appConfig.Port))
}

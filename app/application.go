package app

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

func Application(opts ...ConfigOption) *application {
    for _, opt := range opts {
        opt(config)
    }
    prepareConfig(config)
    prepareDB(config)

    app := iris.New()
    app.Logger().SetLevel(config.LogLevel)
    app.Configure(iris.WithSocketSharding,
        iris.WithoutBodyConsumptionOnUnmarshal)
    app.UseRouter(requestid.New())
    app.UseRouter(recover.New())
    if app.Logger().Level == golog.DebugLevel {
        app.UseRouter(logger.New())
    }
    app.Get("/", func(ctx iris.Context) {
        // HTTP服务怎么就被nmap探测到是Golang实现的呢？
        // https://github.com/bingoohuang/blog/issues/174
    })
    mvc.Configure(app.Party(config.ContextPath),
        dependencyConfigurator,
        middlewareConfigurator,
        controllerConfigurator)

    return &application{app: app}
}

func (a *application) App() *iris.Application {
    return a.app
}

func (a *application) Run() {
    _ = a.app.Listen(":" + gokits.StrFromInt(config.Port))
}

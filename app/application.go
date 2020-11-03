package app

import (
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/gokits"
    "github.com/kataras/golog"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/middleware/logger"
    "github.com/kataras/iris/v12/middleware/recover"
    "github.com/kataras/iris/v12/middleware/requestid"
    "github.com/kataras/iris/v12/mvc"
)

type application struct {
    app    *iris.Application
    config *base.Config
}

func Application(opts ...base.ConfigOption) *application {
    config := base.PrepareConfig(opts...)
    db := base.PrepareDB(config)
    components := PrepareComponents(config, db)

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
    mvc.Configure(app.Party(config.ContextPath), components.Configurator)

    return &application{app: app, config: config}
}

func (a *application) App() *iris.Application {
    return a.app
}

func (a *application) Run() error {
    return a.app.Listen(":" + gokits.StrFromInt(a.config.Port))
}

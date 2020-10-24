package lannister

import (
    "flag"
    "github.com/BurntSushi/toml"
    "github.com/CharLemAznable/gokits"
    "github.com/kataras/golog"
)

type AppConfig struct {
    gokits.HttpServerConfig

    LogLevel string

    DriverName             string
    DataSourceName         string
    MaxOpenConns           int
    MaxIdleConns           int
    ConnMaxIdleTimeInMills int64
    ConnMaxLifetimeInMills int64
}

var (
    _configFile string
    appConfig   AppConfig
)

func init() {
    flag.StringVar(&_configFile, "configFile",
        "appConfig.toml", "config file path")
    flag.Parse()

    loadAppConfig(_configFile, &appConfig)
    fixedAppConfig(&appConfig)

    golog.SetLevel(appConfig.LogLevel)
    golog.Infof("appConfig: %+v", appConfig)
}

func loadAppConfig(configFile string, appConfig *AppConfig) {
    if _, err := toml.DecodeFile(configFile, appConfig); err != nil {
        golog.Errorf("config file decode error: %s", err.Error())
    }
}

func fixedAppConfig(appConfig *AppConfig) {
    gokits.If(0 == appConfig.Port, func() {
        appConfig.Port = 4791
    })
    gokits.If("" == appConfig.ContextPath, func() {
        appConfig.ContextPath = "lannister"
    })
    gokits.If("" == appConfig.LogLevel, func() {
        appConfig.LogLevel = "info"
    })
}

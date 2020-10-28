package app

import (
    "flag"
    "github.com/BurntSushi/toml"
    "github.com/CharLemAznable/gokits"
    "github.com/kataras/golog"
    "testing"
)

type (
    Config struct {
        gokits.HttpServerConfig

        LogLevel string

        DriverName             string
        DataSourceName         string
        MaxOpenConns           int
        MaxIdleConns           int
        ConnMaxIdleTimeInMills int64
        ConnMaxLifetimeInMills int64
    }

    ConfigOption func(*Config)
)

var config = &Config{}

func init() {
    configFile := ""
    flag.StringVar(&configFile, "configFile",
        "config.toml", "config file path")
    flag.Parse()
    if _, err := toml.DecodeFile(configFile, config); err != nil {
        golog.Errorf("config file decode error: %s", err.Error())
    }
}

func prepareConfig(config *Config) {
    fixedConfig(config)

    golog.SetLevel(config.LogLevel)
    golog.Infof("config: %+v", *config)
}

func fixedConfig(config *Config) {
    gokits.If(0 == config.Port, func() {
        config.Port = 4791
    })
    gokits.If("" == config.ContextPath, func() {
        config.ContextPath = "lannister"
    })
    gokits.If("" == config.LogLevel, func() {
        config.LogLevel = "info"
    })
}

var _ = func() bool {
    testing.Init()
    return true
}()

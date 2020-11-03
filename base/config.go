package base

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

        AccessorVerifyCacheInMills         int64
        MerchantVerifyCacheInMills         int64
        AccessorMerchantVerifyCacheInMills int64
    }

    ConfigOption func(*Config)
)

var globalConfig = &Config{}

func init() {
    testing.Init()
    configFile := ""
    flag.StringVar(&configFile, "configFile",
        "config.toml", "config file path")
    flag.Parse()
    if _, err := toml.DecodeFile(configFile, globalConfig); err != nil {
        golog.Errorf("config file decode error: %s", err.Error())
    }
}

func PrepareConfig(opts ...ConfigOption) *Config {
    config := new(Config)
    *config = *globalConfig // 值拷贝

    for _, opt := range opts {
        opt(config)
    }

    fixedConfig(config)

    golog.SetLevel(config.LogLevel)
    golog.Infof("config: %+v", *config)

    return config
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

    // 默认缓存1min
    gokits.If(0 == config.AccessorVerifyCacheInMills, func() {
        config.AccessorVerifyCacheInMills = 60 * 1000
    })
    gokits.If(0 == config.MerchantVerifyCacheInMills, func() {
        config.MerchantVerifyCacheInMills = 60 * 1000
    })
    gokits.If(0 == config.AccessorMerchantVerifyCacheInMills, func() {
        config.AccessorMerchantVerifyCacheInMills = 60 * 1000
    })
}

package base

import (
    "flag"
    "github.com/BurntSushi/toml"
    "github.com/CharLemAznable/gokits"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/golog"
    "testing"
    "time"
)

/*************** beans          ***************/

type BaseResp struct {
    ErrorCode string `json:"errorCode,omitempty"`
    ErrorDesc string `json:"errorDesc,omitempty"`
}

/*************** config         ***************/

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

/*************** database       ***************/

func PrepareDB(config *Config) *sqlx.DB {
    db := loadSqlxDB(config)
    if nil == db {
        db = sqlx.NewDb(nil, "")
    }

    db.MapperFunc(func(s string) string { return s })
    return db
}

func loadSqlxDB(config *Config) *sqlx.DB {
    db, err := sqlx.Open(config.DriverName, config.DataSourceName)
    if err != nil {
        golog.Errorf("open sqlx.DB error: %s", err.Error())
        return nil
    }

    db.SetMaxOpenConns(config.MaxOpenConns)
    db.SetMaxIdleConns(config.MaxIdleConns)
    db.SetConnMaxIdleTime(time.Millisecond * time.Duration(config.ConnMaxIdleTimeInMills))
    db.SetConnMaxLifetime(time.Millisecond * time.Duration(config.ConnMaxLifetimeInMills))

    if err = db.Ping(); err != nil {
        golog.Errorf("connect DB error: %s", err.Error())
        return nil
    }
    golog.Infof("DB: %+v", db)
    return db
}

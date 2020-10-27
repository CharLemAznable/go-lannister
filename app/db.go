package app

import (
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/golog"
    "time"
)

func prepareDB(config *Config) {
    db := LoadSqlxDB(config)
    if nil == db {
        db = sqlx.NewDb(nil, "")
    }

    db.MapperFunc(func(s string) string { return s })
    RegisterDependency("db", db)
}

func LoadSqlxDB(config *Config) *sqlx.DB {
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

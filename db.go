package lannister

import (
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/golog"
    "time"
)

func preparingDB() {
    db := loadSqlxDB(&appConfig)
    if nil != db {
        RegisterDependency("db", db)
    } else {
        RegisterDependency("db", sqlx.NewDb(nil, "dummy"))
    }
}

func loadSqlxDB(appConfig *AppConfig) *sqlx.DB {
    db, err := sqlx.Open(appConfig.DriverName, appConfig.DataSourceName)
    if err != nil {
        golog.Errorf("open sqlx.DB error: %s", err.Error())
        return nil
    }

    db.SetMaxOpenConns(appConfig.MaxOpenConns)
    db.SetMaxIdleConns(appConfig.MaxIdleConns)
    db.SetConnMaxIdleTime(time.Millisecond * time.Duration(appConfig.ConnMaxIdleTimeInMills))
    db.SetConnMaxLifetime(time.Millisecond * time.Duration(appConfig.ConnMaxLifetimeInMills))

    if err = db.Ping(); err != nil {
        golog.Errorf("connect DB error: %s", err.Error())
        return nil
    }
    golog.Infof("DB: %+v", db)
    return db
}

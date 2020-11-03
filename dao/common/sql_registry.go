package common

import (
    "errors"
    "github.com/CharLemAznable/go-lannister/base"
    "github.com/CharLemAznable/sqlx"
)

type SqlRegistry struct {
    *base.Registry
}

func NewSqlRegistry(name string) *SqlRegistry {
    return &SqlRegistry{base.NewRegistry(name)}
}

func (r *SqlRegistry) GetSql(db *sqlx.DB) interface{} {
    if nil == db {
        panic(errors.New("sqlx.DB is nil"))
    }
    driverName := db.DriverName()
    return r.Get(driverName)
}

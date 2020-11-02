package elf

import (
    "errors"
    "github.com/CharLemAznable/sqlx"
)

type SqlBundleRegistry struct {
    *Registry
}

func NewSqlBundleRegistry(name string) *SqlBundleRegistry {
    return &SqlBundleRegistry{NewRegistry(name)}
}

func (r *SqlBundleRegistry) GetSqlBundle(db *sqlx.DB) interface{} {
    if nil == db {
        panic(errors.New("sqlx.DB is nil"))
    }
    driverName := db.DriverName()
    return r.Get(driverName)
}

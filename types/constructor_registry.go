package types

import (
    "errors"
    "fmt"
    "github.com/CharLemAznable/sqlx"
)

type DaoConstructorRegistry struct {
    *Registry
}

func NewDaoConstructorRegistry(name string) *DaoConstructorRegistry {
    return &DaoConstructorRegistry{NewRegistry(name)}
}

func (r *DaoConstructorRegistry) GetDaoConstructor(db *sqlx.DB) interface{} {
    if nil == db {
        panic(errors.New("sqlx.DB is nil"))
    }
    driverName := db.DriverName()
    constructor := r.Get(driverName)
    if nil == constructor {
        panic(errors.New(fmt.Sprintf("Unknown %s for driver %q "+
            "(forgotten import?)", r.name, driverName)))
    }
    return constructor
}

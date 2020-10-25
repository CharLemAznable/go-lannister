package types

import (
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/golog"
)

type DaoConstructorRegistry struct {
    *Registry
}

func NewDaoConstructorRegistry(name string) *DaoConstructorRegistry {
    return &DaoConstructorRegistry{NewRegistry(name)}
}

func (r *DaoConstructorRegistry) GetDaoConstructor(db *sqlx.DB) interface{} {
    if nil == db {
        golog.Error("Nil sqlx.DB")
        return nil
    }
    driverName := db.DriverName()
    constructor := r.Get(driverName)
    if nil == constructor {
        golog.Errorf("Unknown %s for driver %q "+
            "(forgotten import?)", r.name, driverName)
        return nil
    }
    return constructor
}

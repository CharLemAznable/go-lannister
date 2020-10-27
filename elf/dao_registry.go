package elf

import (
    "errors"
    "fmt"
    "github.com/CharLemAznable/sqlx"
    "reflect"
)

type DaoRegistry struct {
    *Registry
}

func NewDaoRegistry(name string) *DaoRegistry {
    return &DaoRegistry{NewRegistry(name)}
}

func (r *DaoRegistry) GetDao(db *sqlx.DB) interface{} {
    if nil == db {
        panic(errors.New("sqlx.DB is nil"))
    }
    driverName := db.DriverName()
    builder := r.Get(driverName)
    if nil == builder {
        panic(errors.New(fmt.Sprintf("Unknown %s for driver %q "+
            "(forgotten import?)", r.name, driverName)))
    }

    v := reflect.ValueOf(builder)
    t := v.Type()
    if reflect.Func != v.Kind() || 1 != t.NumIn() || 1 != t.NumOut() {
        panic(errors.New("registered item should be: func(*sqlx.DB) DaoType"))
    }

    inputs := []reflect.Value{reflect.ValueOf(db)}
    outputs := v.Call(inputs)
    return outputs[0].Interface()
}

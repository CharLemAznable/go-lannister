package base

import (
    "errors"
    "fmt"
    "github.com/CharLemAznable/sqlx"
    "github.com/kataras/golog"
    "reflect"
    "sync"
)

type Registry struct {
    sync.RWMutex

    name  string
    table map[string]interface{}
}

func NewRegistry(name string) *Registry {
    return &Registry{name: name,
        table: make(map[string]interface{})}
}

func (r *Registry) Register(name string, item interface{}) {
    r.Lock()
    defer r.Unlock()

    if item == nil {
        golog.Errorf("Register %s is nil", r.name)
        return
    }
    if _, dup := r.table[name]; dup {
        golog.Errorf("Register %s duplicated for %s", r.name, name)
        return
    }
    r.table[name] = item
}

func (r *Registry) Get(name string) interface{} {
    r.RLock()
    defer r.RUnlock()
    return r.table[name]
}

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

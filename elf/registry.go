package elf

import (
    "github.com/kataras/golog"
    "sort"
    "sync"
)

type Registry struct {
    sync.RWMutex

    name  string
    table map[string]interface{}
}

type Consumer func(interface{})

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

func (r *Registry) RegisterCover(name string, item interface{}) {
    r.Lock()
    defer r.Unlock()
    r.table[name] = item
}

func (r *Registry) Get(name string) interface{} {
    r.RLock()
    defer r.RUnlock()
    return r.table[name]
}

func (r *Registry) Iterate(consumer Consumer) {
    if nil == consumer {
        return
    }
    r.RLock()
    defer r.RUnlock()

    for _, item := range r.table {
        consumer(item)
    }
}

func (r *Registry) IterateSorted(consumer Consumer) {
    if nil == consumer {
        return
    }
    r.RLock()
    defer r.RUnlock()

    names := make([]string, 0)
    for name := range r.table {
        names = append(names, name)
    }
    sort.Strings(names)

    for _, name := range names {
        consumer(r.table[name])
    }
}

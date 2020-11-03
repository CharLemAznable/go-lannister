package base

import (
    "github.com/CharLemAznable/sqlx"
    "github.com/stretchr/testify/assert"
    "testing"
)

type (
    TestComponent1 struct{}
    TestComponent2 struct{}
)

var registry = NewRegistry("Component")

func TestRegistry(t *testing.T) {
    a := assert.New(t)

    registry.Register("nil", nil)
    a.Nil(registry.Get("nil"))

    component1 := &TestComponent1{}
    component2 := &TestComponent2{}
    registry.Register("same", component1)
    registry.Register("same", component2)
    a.Equal(component1, registry.Get("same"))
}

var daoRegistry = NewDaoRegistry("Dao")

func TestDaoRegistry(t *testing.T) {
    a := assert.New(t)

    a.Panics(func() {
        daoRegistry.GetDao(nil)
    })

    db := sqlx.NewDb(nil, "foo")
    a.Panics(func() {
        daoRegistry.GetDao(db)
    })

    daoRegistry.Register("foo", "error")
    a.Panics(func() {
        daoRegistry.GetDao(db)
    })

    db = sqlx.NewDb(nil, "bar")
    flag := ""
    daoRegistry.Register("bar",
        func(db *sqlx.DB) interface{} {
            flag = "bar"
            return nil
        })
    a.Nil(daoRegistry.GetDao(db))
    a.Equal("bar", flag)
}

package elf

import (
    "github.com/CharLemAznable/sqlx"
    "github.com/stretchr/testify/assert"
    "testing"
)

var daoRegistry = NewDaoRegistry("Dao")

func TestDaoConstructorRegistry(t *testing.T) {
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
